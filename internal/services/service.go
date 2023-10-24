package services

import (
	"FioapiKafka/internal/models"
	"FioapiKafka/internal/repository"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type PersonService interface {
	GetPeople() ([]*models.PersonOut, error)
	GetPersonByID(id string) (*models.PersonOut, error)
	GetPersonsByName(name string) ([]*models.PersonOut, error)
	GetPersonsByAge(age string) ([]*models.PersonOut, error)
	GetPersonsByGender(gender string) ([]*models.PersonOut, error)
	UpdatePerson(id string, person *models.PersonOut) error
	DeletePerson(id string) error
	CreatePeople(aproovedFIO chan *models.PersonOut)
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

func (s *AddService) GetPeople() ([]*models.PersonOut, error) {
	return s.repo.GetPeople()
}

func (s *AddService) GetPersonByID(id string) (*models.PersonOut, error) {
	return s.repo.GetPersonByID(id)
}

func (s *AddService) GetPersonsByName(name string) ([]*models.PersonOut, error) {
	return s.repo.GetPersonsByName(name)
}

func (s *AddService) GetPersonsByAge(age string) ([]*models.PersonOut, error) {
	return s.repo.GetPersonsByAge(age)
}

func (s *AddService) GetPersonsByGender(gender string) ([]*models.PersonOut, error) {
	return s.repo.GetPersonsByGender(gender)
}

func (s *AddService) UpdatePerson(id string, person *models.PersonOut) error {
	return s.repo.UpdatePerson(id, person)
}

func (s *AddService) DeletePerson(id string) error {
	return s.repo.DeletePerson(id)
}

func (s *AddService) CreatePeople(aproovedFIO chan *models.PersonOut) {
	for data := range aproovedFIO {
		log.Infoln("A new message was received from the channel aproovedFIO")
		s.repo.CreatePeople(data)
	}
}
