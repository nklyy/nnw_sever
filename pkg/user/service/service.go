package service

import (
	"nnw_s/config"
	"nnw_s/pkg/user/model"
	"nnw_s/pkg/user/repository"
)

type User interface {
	GetUserById(userId string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetTemplateUserDataById(uid string) (*model.TemplateData, error)

	CreateUser(email string, password string, OTPKey string) (*string, error)
	CreateTemplateUserData(secret string) (*string, error)
}

type Service struct {
	User
}

func NewService(repos *repository.Repository, cfg config.Configurations) *Service {
	return &Service{
		User: NewUserService(repos.User, cfg),
	}
}
