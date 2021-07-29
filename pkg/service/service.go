package service

import (
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
)

type User interface {
	GetUserById(userId string) (*model.User, error)
	CreateUser(user model.User) (*string, error)
	UpdateUser(user model.User) error
}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
