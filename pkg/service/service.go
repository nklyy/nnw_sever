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
	GetUserByEmail(email string) (*model.User, error)
	GetTemplateUserDataById(uid string) (*model.TemplateData, error)

	CreateUser(email string, password string, OTPKey string) (*string, error)
	CreateTemplateUserData(secret string) (*string, error)

	CreateJWTToken(email string) (string, error)
	VerifyJWTToken(id string) (*string, error)

	Generate2FaImage(email string) (*bytes.Buffer, *otp.Key, error)
	Check2FaCode(code string, secret string) bool

	CheckPassword(password string, hashPassword string) (bool, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, cfg config.Configurations, emailClient config.SMTPClient) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, cfg, emailClient),
	}
}
