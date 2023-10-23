package repository

import (
	"FioapiKafka/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type repository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAddDb(db *sql.DB, rdb *redis.Client) PersonRepository {
	return &repository{db: db, rdb: rdb}
}

func (r *repository) CreatePeople(data *models.PersonOut) {
	log.Infoln("Adding data to database")
	Natio, err := json.Marshal(data.Nationality)
	if err != nil {
		log.Errorf("error marshal data.Nationality:%s", err.Error())
		return
	}
	var id string
	err = r.db.QueryRow("INSERT INTO fioapi (name,surname,patronymic,age,gender,nationality) values ($1,$2,$3,$4,$5,$6) RETURNING id",
		data.Name, data.Surname, data.Patronymic, data.Age, data.Gender, Natio).Scan(&id)
	if err != nil {
		log.Errorf("error insert data in db:%s", err.Error())
		return
	}
	log.Infof("Adding data to redis, key:%s", id)
	dataTOredis, err := json.Marshal(data)
	err = r.rdb.Set(ctx, id, dataTOredis, 0).Err()
	if err != nil {
		log.Errorf("error setting data in redis:%s", err.Error())
		return
	}
}

func (r *repository) GetPeople() ([]*models.PersonOut, error) {
	log.Infoln("Get all persons from cache")
	keys, err := r.rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Errorf("error getting all data from redis:%s", err.Error())
		return nil, err
	}
	var AllData []*models.PersonOut
	for _, key := range keys {
		value, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			log.Errorf("error retrieving data by key %s:%s", key, err.Error())
			continue
		}
		data := new(models.PersonOut)
		err = json.Unmarshal([]byte(value), data)
		if err != nil {
			log.Errorf("error unmarshal data to *models.PersonOut:%s", err.Error())
			continue
		}
		AllData = append(AllData, data)
	}
	return AllData, nil
}

func (r *repository) GetPersonByID(id string) (*models.PersonOut, error) {
	log.Infof("Geting data from cache by id:%s", id)
	value, err := r.rdb.Get(ctx, id).Result()
	if err != nil {
		log.Errorf("error geting data from cache by id - %s:%s", id, err.Error())
		return nil, err
	}
	data := new(models.PersonOut)
	err = json.Unmarshal([]byte(value), data)
	if err != nil {
		log.Errorf("error unmarshal data to *models.PersonOut:%s", err.Error())
	}
	return data, nil
}

func (r *repository) GetPersonsByName(name string) ([]*models.PersonOut, error) {
	log.Infof("Get persons by name from database, name - %s", name)
	rows, err := r.db.Query("SELECT * FROM fioapi WHERE name = $1", name)
	if err != nil {
		log.Errorf("err selected in database:%s", err.Error())
		return nil, err
	}
	defer rows.Close()
	persons := make([]*models.PersonOut, 0)
	for rows.Next() {
		value := &models.PersonOut{}
		var Natio []byte
		var id int
		err := rows.Scan(&id, &value.Name, &value.Surname, &value.Patronymic, &value.Age, &value.Gender, &Natio)
		if err != nil {
			log.Errorf("error scan rows:%s", err.Error())
			continue
		}
		err = json.Unmarshal(Natio, &value.Nationality)
		if err != nil {
			log.Errorf("error unmarshal type - Nationality:%s", err.Error())
			continue
		}
		persons = append(persons, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *repository) GetPersonsByAge(age string) ([]*models.PersonOut, error) {
	log.Infof("Get persons by age from database, age - %s", age)
	rows, err := r.db.Query("SELECT * FROM fioapi WHERE age = $1", age)
	if err != nil {
		log.Errorf("error selected from database:%s", err.Error())
		return nil, err
	}
	defer rows.Close()
	persons := make([]*models.PersonOut, 0)
	for rows.Next() {
		value := &models.PersonOut{}
		var Natio []byte
		var id int
		err := rows.Scan(&id, &value.Name, &value.Surname, &value.Patronymic, &value.Age, &value.Gender, &Natio)
		if err != nil {
			log.Errorf("error scan rows:%s", err.Error())
			continue
		}
		err = json.Unmarshal(Natio, &value.Nationality)
		if err != nil {
			log.Errorf("error unmarshal type - Nationality:%s", err.Error())
			continue
		}
		persons = append(persons, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *repository) GetPersonsByGender(gender string) ([]*models.PersonOut, error) {
	log.Infof("Get persons by gender from database, gender - %s", gender)
	rows, err := r.db.Query("SELECT * FROM fioapi WHERE gender = $1", gender)
	if err != nil {
		log.Errorf("err selected in database:%s", err.Error())
		return nil, err
	}
	defer rows.Close()
	persons := make([]*models.PersonOut, 0)
	for rows.Next() {
		value := &models.PersonOut{}
		var Natio []byte
		var id int
		err := rows.Scan(&id, &value.Name, &value.Surname, &value.Patronymic, &value.Age, &value.Gender, &Natio)
		if err != nil {
			log.Errorf("error scan rows:%s", err.Error())
			continue
		}
		err = json.Unmarshal(Natio, &value.Nationality)
		if err != nil {
			log.Errorf("error unmarshal type - Nationality:%s", err.Error())
			continue
		}
		persons = append(persons, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *repository) UpdatePerson(id string, person *models.PersonOut) error {
	log.Infof("Start update person by id:%s", id)

	NationJSON, err := json.Marshal(person.Nationality)
	if err != nil {
		log.Errorf("error marshaling nationality:%s", err.Error())
		return err
	}

	_, err = r.db.Exec("UPDATE fioapi SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7",
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, NationJSON, id)
	if err != nil {
		log.Errorf("error updating data in db:%s", err.Error())
		return err
	}

	log.Infof("Successfully updated in the database, now updating in the redis")

	dataToRedis, err := json.Marshal(person)
	if err != nil {
		log.Errorf("error marshaling person:%s", err.Error())
		return err
	}

	err = r.rdb.Set(ctx, id, dataToRedis, 0).Err()
	if err != nil {
		log.Errorf("error setting data in redis: %s", err.Error())
		return err
	}

	log.Infof("Successfully updated person with id: %s", id)

	return nil
}

func (r *repository) DeletePerson(id string) error {
	log.Infof("Start delete person by id:%s", id)

	// Delete from db
	_, err := r.db.Exec("DELETE FROM fioapi WHERE id = $1", id)
	if err != nil {
		log.Errorf("error deleting data from db:%s", err.Error())
		return err
	}

	log.Infof("Successfully deleted from the database, now deleting from redis")

	// Delete from Redis
	err = r.rdb.Del(ctx, id).Err()
	if err != nil {
		log.Errorf("error deleting data from Redis:%s", err.Error())
		return err
	}

	log.Infof("Successfully deleted person with id: %s", id)

	return nil
}
