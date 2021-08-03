package service

import (
	"bytes"
	"github.com/pquerna/otp"
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
)

type Authorization interface {
	GetUserById(userId string) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetTemplateUserDataById(uid string) (*model.TemplateData, error)

	CreateUser(login string, email string, password string, OTPKey string) (*string, error)
	CreateTemplateUserData(secret string) (*string, error)

	Generate2FaImage() (*bytes.Buffer, *otp.Key, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
