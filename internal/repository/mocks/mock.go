// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	models "FioapiKafka/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersonRepository is a mock of PersonRepository interface.
type MockPersonRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPersonRepositoryMockRecorder
}

// MockPersonRepositoryMockRecorder is the mock recorder for MockPersonRepository.
type MockPersonRepositoryMockRecorder struct {
	mock *MockPersonRepository
}

// NewMockPersonRepository creates a new mock instance.
func NewMockPersonRepository(ctrl *gomock.Controller) *MockPersonRepository {
	mock := &MockPersonRepository{ctrl: ctrl}
	mock.recorder = &MockPersonRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonRepository) EXPECT() *MockPersonRepositoryMockRecorder {
	return m.recorder
}

// CreatePeople mocks base method.
func (m *MockPersonRepository) CreatePeople(data *models.PersonOut) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePeople", data)
}

// CreatePeople indicates an expected call of CreatePeople.
func (mr *MockPersonRepositoryMockRecorder) CreatePeople(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePeople", reflect.TypeOf((*MockPersonRepository)(nil).CreatePeople), data)
}

// DeletePerson mocks base method.
func (m *MockPersonRepository) DeletePerson(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePerson", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePerson indicates an expected call of DeletePerson.
func (mr *MockPersonRepositoryMockRecorder) DeletePerson(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePerson", reflect.TypeOf((*MockPersonRepository)(nil).DeletePerson), id)
}

// GetPeople mocks base method.
func (m *MockPersonRepository) GetPeople() ([]*models.PersonOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeople")
	ret0, _ := ret[0].([]*models.PersonOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeople indicates an expected call of GetPeople.
func (mr *MockPersonRepositoryMockRecorder) GetPeople() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeople", reflect.TypeOf((*MockPersonRepository)(nil).GetPeople))
}

// GetPersonByID mocks base method.
func (m *MockPersonRepository) GetPersonByID(id string) (*models.PersonOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonByID", id)
	ret0, _ := ret[0].(*models.PersonOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonByID indicates an expected call of GetPersonByID.
func (mr *MockPersonRepositoryMockRecorder) GetPersonByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonByID", reflect.TypeOf((*MockPersonRepository)(nil).GetPersonByID), id)
}

// GetPersonsByAge mocks base method.
func (m *MockPersonRepository) GetPersonsByAge(age string) ([]*models.PersonOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonsByAge", age)
	ret0, _ := ret[0].([]*models.PersonOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonsByAge indicates an expected call of GetPersonsByAge.
func (mr *MockPersonRepositoryMockRecorder) GetPersonsByAge(age interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonsByAge", reflect.TypeOf((*MockPersonRepository)(nil).GetPersonsByAge), age)
}

// GetPersonsByGender mocks base method.
func (m *MockPersonRepository) GetPersonsByGender(gender string) ([]*models.PersonOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonsByGender", gender)
	ret0, _ := ret[0].([]*models.PersonOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonsByGender indicates an expected call of GetPersonsByGender.
func (mr *MockPersonRepositoryMockRecorder) GetPersonsByGender(gender interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonsByGender", reflect.TypeOf((*MockPersonRepository)(nil).GetPersonsByGender), gender)
}

// GetPersonsByName mocks base method.
func (m *MockPersonRepository) GetPersonsByName(name string) ([]*models.PersonOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonsByName", name)
	ret0, _ := ret[0].([]*models.PersonOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonsByName indicates an expected call of GetPersonsByName.
func (mr *MockPersonRepositoryMockRecorder) GetPersonsByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonsByName", reflect.TypeOf((*MockPersonRepository)(nil).GetPersonsByName), name)
}

// UpdatePerson mocks base method.
func (m *MockPersonRepository) UpdatePerson(id string, person *models.PersonOut) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePerson", id, person)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePerson indicates an expected call of UpdatePerson.
func (mr *MockPersonRepositoryMockRecorder) UpdatePerson(id, person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePerson", reflect.TypeOf((*MockPersonRepository)(nil).UpdatePerson), id, person)
}
