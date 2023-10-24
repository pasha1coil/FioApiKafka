package repository

import (
	"FioapiKafka/internal/models"
	"database/sql"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type PersonRepository interface {
	CreatePeople(data *models.PersonOut)
	GetPeople() ([]*models.PersonOut, error)
	GetPersonByID(id string) (*models.PersonOut, error)
	GetPersonsByName(name string) ([]*models.PersonOut, error)
	GetPersonsByAge(age string) ([]*models.PersonOut, error)
	GetPersonsByGender(gender string) ([]*models.PersonOut, error)
	UpdatePerson(id string, person *models.PersonOut) error
	DeletePerson(id string) error
}

type Repository struct {
	PersonRepository
}

func NewRepository(db *sql.DB, rdb *redis.Client) *Repository {
	return &Repository{
		PersonRepository: NewAddDb(db, rdb),
	}
}
