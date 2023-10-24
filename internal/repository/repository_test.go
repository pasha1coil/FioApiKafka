package repository

import (
	"FioapiKafka/internal/models"
	mock_repository "FioapiKafka/internal/repository/mocks"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePeople(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)
	mockPerson := &models.PersonOut{Name: "Test Name", Surname: "Test Surname"}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().CreatePeople(mockPerson).Times(1)

		repo := Repository{PersonRepository: mockRepo}

		repo.CreatePeople(mockPerson)
	})
}

func TestGetPersonByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)

	testCases := []struct {
		name           string
		redisValue     *models.PersonOut
		redisReturnErr error
		expectedPerson *models.PersonOut
		expectedErr    bool
	}{
		{
			name:           "Good",
			redisValue:     &models.PersonOut{Name: "Test", Age: 30},
			redisReturnErr: nil,
			expectedPerson: &models.PersonOut{Name: "Test1", Age: 31},
			expectedErr:    false,
		},
		{
			name:           "Bad",
			redisValue:     nil,
			redisReturnErr: redis.Nil,
			expectedPerson: nil,
			expectedErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			id := "1"

			mockRepo.EXPECT().GetPersonByID(id).Return(tc.redisValue, tc.redisReturnErr)

			_, actualErr := mockRepo.GetPersonByID(id)

			if tc.expectedErr {
				require.Error(t, actualErr)
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}

func TestGetPersonsByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)

	testCases := []struct {
		name            string
		dbValue         []*models.PersonOut
		dbReturnErr     error
		expectedPersons []*models.PersonOut
		expectedErr     bool
	}{
		{
			name:            "Good",
			dbValue:         []*models.PersonOut{{Name: "Test", Age: 30}},
			dbReturnErr:     nil,
			expectedPersons: []*models.PersonOut{{Name: "Test", Age: 30}},
			expectedErr:     false,
		},
		{
			name:            "Bad",
			dbValue:         nil,
			dbReturnErr:     sql.ErrNoRows,
			expectedPersons: nil,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nameCriteria := "Test"

			mockRepo.EXPECT().GetPersonsByName(nameCriteria).Return(tc.dbValue, tc.dbReturnErr)

			actualPersons, actualErr := mockRepo.GetPersonsByName(nameCriteria)

			if tc.expectedErr {
				require.Error(t, actualErr)
			} else {
				require.NoError(t, actualErr)
				require.Equal(t, tc.expectedPersons, actualPersons)
			}
		})
	}
}

func TestGetPersonsByAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)

	testCases := []struct {
		name            string
		dbValue         []*models.PersonOut
		dbReturnErr     error
		expectedPersons []*models.PersonOut
		expectedErr     bool
	}{
		{
			name:            "Good",
			dbValue:         []*models.PersonOut{{Name: "Test", Age: 30}},
			dbReturnErr:     nil,
			expectedPersons: []*models.PersonOut{{Name: "Test", Age: 30}},
			expectedErr:     false,
		},
		{
			name:            "Bad",
			dbValue:         nil,
			dbReturnErr:     sql.ErrNoRows,
			expectedPersons: nil,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ageCriteria := "30"

			mockRepo.EXPECT().GetPersonsByAge(ageCriteria).Return(tc.dbValue, tc.dbReturnErr)

			actualPersons, actualErr := mockRepo.GetPersonsByAge(ageCriteria)

			if tc.expectedErr {
				require.Error(t, actualErr)
			} else {
				require.NoError(t, actualErr)
				require.Equal(t, tc.expectedPersons, actualPersons)
			}
		})
	}
}

func TestGetPersonsByGender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)

	testCases := []struct {
		name            string
		dbValue         []*models.PersonOut
		dbReturnErr     error
		expectedPersons []*models.PersonOut
		expectedErr     bool
	}{
		{
			name:            "Good",
			dbValue:         []*models.PersonOut{{Name: "Test", Gender: "male"}},
			dbReturnErr:     nil,
			expectedPersons: []*models.PersonOut{{Name: "Test", Gender: "male"}},
			expectedErr:     false,
		},
		{
			name:            "Bad",
			dbValue:         nil,
			dbReturnErr:     sql.ErrNoRows,
			expectedPersons: nil,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genderCriteria := "male"

			mockRepo.EXPECT().GetPersonsByGender(genderCriteria).Return(tc.dbValue, tc.dbReturnErr)

			actualPersons, actualErr := mockRepo.GetPersonsByGender(genderCriteria)

			if tc.expectedErr {
				require.Error(t, actualErr)
			} else {
				require.NoError(t, actualErr)
				require.Equal(t, tc.expectedPersons, actualPersons)
			}
		})
	}
}

func TestDeletePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)
	repo := Repository{
		PersonRepository: mockRepo,
	}

	testCases := []struct {
		name      string
		testID    string
		mockErr   error
		expectErr bool
	}{
		{name: "Good", testID: "test1", mockErr: nil, expectErr: false},
		{name: "Bad", testID: "test1", mockErr: fmt.Errorf("Error deleting"), expectErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().DeletePerson(tc.testID).Return(tc.mockErr)

			err := repo.DeletePerson(tc.testID)

			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, tc.mockErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPersonRepository(ctrl)
	repo := Repository{
		PersonRepository: mockRepo,
	}

	testCases := []struct {
		name           string
		testID         string
		personToUpdate *models.PersonOut
		mockErr        error
		expectErr      bool
	}{
		{
			name:           "Good",
			testID:         "1",
			personToUpdate: &models.PersonOut{Name: "Test", Age: 35},
			mockErr:        nil,
			expectErr:      false,
		},
		{
			name:           "Bad",
			testID:         "1",
			personToUpdate: &models.PersonOut{Name: "Test", Age: 35},
			mockErr:        fmt.Errorf("Error updating"),
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().UpdatePerson(tc.testID, tc.personToUpdate).Return(tc.mockErr)

			err := repo.UpdatePerson(tc.testID, tc.personToUpdate)

			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, tc.mockErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
