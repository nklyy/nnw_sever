package service

import (
	"bytes"
	"github.com/pquerna/otp"
	"nnw_s/config"
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GetUserById(userId string) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetTemplateUserDataById(uid string) (*model.TemplateData, error)

	CreateUser(login string, email string, password string, OTPKey string) (*string, error)
	CreateTemplateUserData(secret string) (*string, error)

	CreateJWTToken(login string) (string, error)
	VerifyJWTToken(id string) (*string, error)

	Generate2FaImage(login string) (*bytes.Buffer, *otp.Key, error)

	CheckPassword(password string, hashPassword string) (bool, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, cfg config.Configurations) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, cfg),
	}
}
