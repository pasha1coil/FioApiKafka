package services

import (
	"FioapiKafka/internal/models"
	"FioapiKafka/internal/repository"
)

type PersonService interface {
	GetPeople() ([]*models.Person, error)
	GetPersonByID(id string) (*models.Person, error)
	UpdatePerson(id string, person *models.Person)
	DeletePerson(id string) error
}

type Service struct {
	PersonService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		PersonService: NewAddService(repos),
	}
}

type AddService struct {
	repo repository.PersonRepository
}

func NewAddService(repo repository.PersonRepository) *AddService {
	return &AddService{repo: repo}
}

func (s *AddService) GetPeople() ([]*models.Person, error) {
	return s.repo.GetPeople()
}
func (s *AddService) GetPersonByID(id string) (*models.Person, error) {
	return s.repo.GetPersonByID(id)
}
func (s *AddService) UpdatePerson(id string, person *models.Person) {
	s.repo.UpdatePerson(id, person)
}
func (s *AddService) DeletePerson(id string) error {
	return s.repo.DeletePerson(id)
}
