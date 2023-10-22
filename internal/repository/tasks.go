package repository

import (
	"FioapiKafka/internal/models"
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAddDb(db *sql.DB, rdb *redis.Client) PersonRepository {
	return &repository{db: db, rdb: rdb}
}

func (r *repository) GetPeople() ([]*models.Person, error) {
	// implement your code here
	return nil, nil
}

func (r *repository) GetPersonByID(id string) (*models.Person, error) {
	// implement your code here
	return nil, nil
}

func (r *repository) UpdatePerson(id string, person *models.Person) {
	// implement your code here
}

func (r *repository) DeletePerson(id string) error {
	// implement your code here
	return nil
}
