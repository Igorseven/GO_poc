package service

import (
	entity "PocGo/internal/domain/entities"
	notify "PocGo/internal/domain/notification"
	repository "PocGo/internal/repositories"
)

const (
	Entity = "Usuário"
)

type UserService interface {
	GetById(id string) (*entity.User, error)
	GetAll(date string) (*[]entity.User, error)
	Update(toUpdate *entity.User) error
	UpdateOldUsersStatus() (int, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}

func (service *userService) GetById(id string) (*entity.User, error) {
	user, err := service.userRepository.FindById(id)

	if err != nil {
		return nil, notify.CreateCustomNotification(notify.NotFound, Entity, err)
	}

	return user, nil
}

func (service *userService) GetAll(date string) (*[]entity.User, error) {
	users, err := service.userRepository.FindAll(date)

	if err != nil {
		return nil, notify.CreateCustomNotification(notify.NotFound, Entity, err)
	}

	return users, nil
}

func (service *userService) Update(dtoUpdate *entity.User) error {
	user, err := service.GetById(dtoUpdate.ID)
	if err != nil {
		return notify.CreateCustomNotification(notify.NotFound, Entity, err)
	}

	if dtoUpdate.Name != "" {
		user.Name = dtoUpdate.Name
	}
	if dtoUpdate.Email != "" {
		user.Email = dtoUpdate.Email
	}

	user.Status = dtoUpdate.Status

	if err := service.userRepository.Update(user); err != nil {
		return notify.CreateCustomNotification(notify.InvalidData, Entity, err)
	}

	*dtoUpdate = *user

	return nil
}

func (service *userService) UpdateOldUsersStatus() (int, error) {
	oldUsers, err := service.userRepository.FindOldUsers()
	if err != nil {
		return 0, notify.CreateCustomNotification(notify.NotFound, Entity, err)
	}

	if oldUsers == nil || len(*oldUsers) == 0 {
		return 0, nil
	}

	updatedCount := 0
	for _, user := range *oldUsers {
		if user.Status == 2 {
			continue
		}

		if err := service.userRepository.UpdateStatus(user.ID, 2); err != nil {
			continue
		}
		updatedCount++
	}

	return updatedCount, nil
}
