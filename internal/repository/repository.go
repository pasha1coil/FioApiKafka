package repository

import (
	"FioapiKafka/internal/models"
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type PersonRepository interface {
	GetPeople() ([]*models.Person, error)
	GetPersonByID(id string) (*models.Person, error)
	UpdatePerson(id string, person *models.Person)
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
