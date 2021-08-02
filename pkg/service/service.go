package service

import (
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
)

type User interface {
	GetUserById(userId string) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetTemplateUserDataById(uid string) (*model.TemplateData, error)
	CreateUser(login string, email string, password string, OTPKey string) (*string, error)
	CreateTemplateUserData(secret string) (*string, error)
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
