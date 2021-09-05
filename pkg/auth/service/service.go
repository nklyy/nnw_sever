package service

import (
	"bytes"
	"github.com/pquerna/otp"
	"nnw_s/config"
	"nnw_s/pkg/auth/repository"
	repository2 "nnw_s/pkg/user/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateJWTToken(email string) (string, error)
	VerifyJWTToken(id string) (*string, error)

	Generate2FaImage(email string) (*bytes.Buffer, *otp.Key, error)
	Check2FaCode(code string, secret string) bool

	CheckPassword(password string, hashPassword string) (bool, error)
}

type Service struct {
	Authorization
}

func NewService(arepos *repository.Repository,
	urepos *repository2.Repository,
	cfg config.Configurations) *Service {
	return &Service{
		Authorization: NewAuthService(arepos.Authorization, urepos.User, cfg),
	}
}
