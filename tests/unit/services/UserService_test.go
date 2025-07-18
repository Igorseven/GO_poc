package services_test

import (
	entity "PocGo/internal/domain/entities"
	notify "PocGo/internal/domain/notification"
	service "PocGo/internal/services"
	"PocGo/tests/helpers"
	"PocGo/tests/mocks"
	"errors"
	"testing"
)

func TestUserService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		date          string
		mockSetup     func(*mocks.UserRepositoryMock)
		expectedUsers *[]entity.User
		expectedError error
	}{
		{
			name: "Success - Users found",
			date: "2025/07/25",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindAllFunc = func(date string) (*[]entity.User, error) {
					return helpers.CreateTestUsers(3), nil
				}
			},
			expectedUsers: helpers.CreateTestUsers(3),
		}, {
			name: "Success - No users found",
			date: "2025/07/25",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindAllFunc = func(date string) (*[]entity.User, error) {
					return helpers.CreateTestUsers(0), nil
				}
			},
			expectedUsers: helpers.CreateTestUsers(0),
		}, {
			name: "Error - Database error",
			date: "2025/07/25",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindAllFunc = func(date string) (*[]entity.User, error) {
					return nil, notify.CreateSimpleNotification(
						notify.FindErrorRepository,
						errors.New("database error"))
				}
			},
			expectedUsers: nil,
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrage
			mockRepo := mocks.NewUserRepositoryMock()
			tt.mockSetup(mockRepo)
			userService := service.NewUserService(mockRepo)

			//Act
			users, err := userService.GetAll(tt.date)

			//Assert
			if tt.expectedError != nil {
				helpers.AssertError(t, err, "Should return an error")
			} else {
				helpers.AssertNotNil(t, users, "Users should not be nil")
				helpers.AssertEqual(t, len(*users), len(*tt.expectedUsers), "User count should match")
			}

			helpers.AssertEqual(t, 1, len(mockRepo.FindAllCalls), "Repository FindAll should be called once")
			helpers.AssertEqual(t, tt.date, mockRepo.FindAllCalls[0], "Repository FindAll should be called with the correct date")
		})
	}
}

func TestUserService_GetById(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		mockSetup     func(*mocks.UserRepositoryMock)
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name:   "Success - User found",
			userID: "1",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindByIdFunc = func(id string) (*entity.User, error) {
					return helpers.CreateTestUser("1"), nil
				}
			},
			expectedUser:  helpers.CreateTestUser("1"),
			expectedError: nil,
		},
		{
			name:   "Error - User not found",
			userID: "999",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindByIdFunc = func(id string) (*entity.User, error) {
					return nil, notify.CreateSimpleNotification(
						notify.NotFound,
						errors.New("user not found"))
				}
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name:   "Error - Database error",
			userID: "1",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindByIdFunc = func(id string) (*entity.User, error) {
					return nil, notify.CreateSimpleNotification(
						notify.FindErrorRepository,
						errors.New("database connection error"))
				}
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewUserRepositoryMock()
			tt.mockSetup(mockRepo)
			userService := service.NewUserService(mockRepo)

			// Act
			user, err := userService.GetById(tt.userID)

			// Assert
			if tt.expectedError != nil {
				helpers.AssertError(t, err, "Should return an error")
			} else {
				helpers.AssertNoError(t, err, "Should not return an error")
				helpers.AssertNotNil(t, user, "User should not be nil")
				helpers.AssertEqual(t, tt.expectedUser.ID, user.ID, "User ID should match")
				helpers.AssertEqual(t, tt.expectedUser.Name, user.Name, "User Name should match")
				helpers.AssertEqual(t, tt.expectedUser.Email, user.Email, "User Email should match")
				helpers.AssertEqual(t, tt.expectedUser.Status, user.Status, "User Status should match")
			}

			helpers.AssertEqual(t, 1, len(mockRepo.FindByIdCalls), "Repository FindById should be called once")
			helpers.AssertEqual(t, tt.userID, mockRepo.FindByIdCalls[0], "Repository FindById should be called with the correct ID")
		})
	}
}

func TestUserService_UpdateOldUsersStatus(t *testing.T) {
	tests := []struct {
		name             string
		mockSetup        func(*mocks.UserRepositoryMock)
		expectedCount    int
		expectedError    error
		expectedStatuses map[string]int
	}{
		{
			name: "Success - Update multiple users",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				users := helpers.CreateTestUsers(3)
				mock.FindOldUsersFunc = func() (*[]entity.User, error) {
					return users, nil
				}
				mock.UpdateStatusFunc = func(id string, status int) error {
					return nil
				}
			},
			expectedCount:    3,
			expectedError:    nil,
			expectedStatuses: map[string]int{"1": 2, "2": 2, "3": 2},
		},
		{
			name: "Success - No users to update",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindOldUsersFunc = func() (*[]entity.User, error) {
					return &[]entity.User{}, nil
				}
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "Partial Success - Some updates fail",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				users := helpers.CreateTestUsers(3)
				mock.FindOldUsersFunc = func() (*[]entity.User, error) {
					return users, nil
				}
				mock.UpdateStatusFunc = func(id string, status int) error {
					if id == "2" {
						return errors.New("update failed")
					}
					return nil
				}
			},
			expectedCount:    2,
			expectedError:    nil,
			expectedStatuses: map[string]int{"1": 2, "2": 2, "3": 2},
		},
		{
			name: "Error - Failed to find old users",
			mockSetup: func(mock *mocks.UserRepositoryMock) {
				mock.FindOldUsersFunc = func() (*[]entity.User, error) {
					return nil, errors.New("database error")
				}
			},
			expectedCount: 0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewUserRepositoryMock()
			tt.mockSetup(mockRepo)
			userService := service.NewUserService(mockRepo)

			// Act
			count, err := userService.UpdateOldUsersStatus()

			// Assert
			if tt.expectedError != nil {
				helpers.AssertError(t, err, "Should return an error")
			} else {
				helpers.AssertNoError(t, err, "Should not return an error")
				helpers.AssertEqual(t, tt.expectedCount, count, "Updated count should match expected")
			}

			helpers.AssertEqual(t, 1, mockRepo.FindOldUsersCalls, "FindOldUsers should be called once")

			if tt.expectedStatuses != nil {
				for id, status := range tt.expectedStatuses {
					if actualStatus, ok := mockRepo.UpdateStatusCalls[id]; ok {
						helpers.AssertEqual(t, status, actualStatus, "Status for user "+id+" should match")
					}
				}
			}
		})
	}
}
