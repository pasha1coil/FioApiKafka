package services

import (
	"FioapiKafka/internal/models"
	mock_services "FioapiKafka/internal/services/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPeople(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTasks := mock_services.NewMockPersonService(ctrl)

	items := []*models.PersonOut{
		{
			Name: "Test1",
		},
		{
			Name: "Test2",
		},
	}

	mockTasks.EXPECT().GetPeople().Return(items, nil).Times(1)

	service := &Service{
		PersonService: mockTasks,
	}

	result, err := service.GetPeople()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != len(items) {
		t.Fatalf("Expected %d items but got %d", len(items), len(result))
	}
}

func TestGetPersonByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTasks := mock_services.NewMockPersonService(ctrl)
	expectedResult := &models.PersonOut{Name: "Test"}
	mockTasks.EXPECT().GetPersonByID("1").Return(expectedResult, nil)
	service := &Service{
		PersonService: mockTasks,
	}
	result, err := service.GetPersonByID("1")
	if err != nil {
		log.Errorf("Error testing get person by id:%s", err.Error())
	}

	assert.Equal(t, expectedResult, result)
}

func TestGetPersonsByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_services.NewMockPersonService(ctrl)
	expectedResult := []*models.PersonOut{
		{Name: "Test"}, {Name: "Test"},
	}
	mockService.EXPECT().GetPersonsByName("Test").Return(expectedResult, nil)
	service := &Service{
		PersonService: mockService,
	}
	result, err := service.GetPersonsByName("Test")
	if err != nil {
		log.Errorf("Error testing get persons by name: %s", err.Error())
	}

	assert.Equal(t, expectedResult, result)
}

func TestGetPersonsByAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_services.NewMockPersonService(ctrl)
	expectedResult := []*models.PersonOut{
		{Name: "Test"}, {Name: "Test2"},
	}
	mockService.EXPECT().GetPersonsByAge("40").Return(expectedResult, nil)
	service := &Service{
		PersonService: mockService,
	}
	result, err := service.GetPersonsByAge("40")
	if err != nil {
		log.Errorf("Error testing get persons by age: %s", err.Error())
	}

	assert.Equal(t, expectedResult, result)
}

func TestGetPersonsByGender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_services.NewMockPersonService(ctrl)
	expectedResult := []*models.PersonOut{
		{Name: "Test"}, {Name: "Test2"},
	}
	mockService.EXPECT().GetPersonsByGender("male").Return(expectedResult, nil)
	service := &Service{
		PersonService: mockService,
	}
	result, err := service.GetPersonsByGender("male")
	if err != nil {
		log.Errorf("Error testing get persons by gender: %s", err.Error())
	}

	assert.Equal(t, expectedResult, result)
}

func TestDeletePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockPersonService(ctrl)

	mockService.EXPECT().DeletePerson("1").Return(nil)

	service := &Service{
		PersonService: mockService,
	}

	err := service.DeletePerson("1")

	if err != nil {
		t.Errorf("Error deleting person: %s", err.Error())
	}
}

func TestUpdatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockPersonService(ctrl)
	person := &models.PersonOut{
		Name: "Test",
	}

	mockService.EXPECT().UpdatePerson("1", person).Return(nil)
	service := &Service{
		PersonService: mockService,
	}

	err := service.UpdatePerson("1", person)
	if err != nil {
		t.Errorf("Error updating person: %s", err.Error())
	}
}
