package service

import (
	repository "PocGo/internal/repositories"
)

type Services struct {
	User UserService
	// Outros serviços aqui
}

func NewServices(repositories *repository.Repositories) *Services {

	return &Services{
		User: NewUserService(repositories.User),
	}
}
