package mocks

import (
	entity "PocGo/internal/domain/entities"
	"errors"
)

type UserRepositoryMock struct {
	FindByIdFunc     func(id string) (*entity.User, error)
	FindAllFunc      func(date string) (*[]entity.User, error)
	UpdateFunc       func(user *entity.User) error
	FindOldUsersFunc func() (*[]entity.User, error)
	UpdateStatusFunc func(id string, status int) error

	FindByIdCalls     []string
	FindAllCalls      []string
	UpdateCalls       []*entity.User
	FindOldUsersCalls int
	UpdateStatusCalls map[string]int
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		FindByIdCalls:     []string{},
		FindAllCalls:      []string{},
		UpdateCalls:       []*entity.User{},
		FindOldUsersCalls: 0,
		UpdateStatusCalls: make(map[string]int),
	}
}

func (mock *UserRepositoryMock) FindById(id string) (*entity.User, error) {
	mock.FindByIdCalls = append(mock.FindByIdCalls, id)
	if mock.FindByIdFunc != nil {
		return mock.FindByIdFunc(id)
	}
	return nil, errors.New("FindByIdFunc not implemented")
}

func (mock *UserRepositoryMock) FindAll(date string) (*[]entity.User, error) {
	mock.FindAllCalls = append(mock.FindAllCalls, date)
	if mock.FindAllFunc != nil {
		return mock.FindAllFunc(date)
	}
	return nil, errors.New("FindAllFunc not implemented")
}

func (mock *UserRepositoryMock) Update(user *entity.User) error {
	mock.UpdateCalls = append(mock.UpdateCalls, user)
	if mock.UpdateFunc != nil {
		return mock.UpdateFunc(user)
	}
	return errors.New("UpdateFunc not implemented")
}

func (mock *UserRepositoryMock) FindOldUsers() (*[]entity.User, error) {
	mock.FindOldUsersCalls++
	if mock.FindOldUsersFunc != nil {
		return mock.FindOldUsersFunc()
	}
	return nil, errors.New("FindOldUsersFunc not implemented")
}

func (mock *UserRepositoryMock) UpdateStatus(id string, status int) error {
	mock.UpdateStatusCalls[id] = status
	if mock.UpdateStatusFunc != nil {
		return mock.UpdateStatusFunc(id, status)
	}
	return errors.New("UpdateStatusFunc not implemented")
}
