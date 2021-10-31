package credentials

import (
	"context"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	CreateCredentials(ctx context.Context, password string, secretOTP SecretOTP) (*DTO, error)
	ValidatePassword(ctx context.Context, credentialsDTO *DTO, password string) error
	DecodePassword(ctx context.Context, password string) (string, error)
}

type service struct {
	log          *logrus.Logger
	shift        int
	passwordSalt int
}

func NewService(log *logrus.Logger, shift, passwordSalt int) (Service, error) {
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	return &service{log: log, shift: shift, passwordSalt: passwordSalt}, nil
}

// CreateCredentials decodes password, hashing it and creates Credentials struct
// There you can also put encrypting logic of secretOTP
func (svc *service) CreateCredentials(ctx context.Context, password string, secretOTP SecretOTP) (*DTO, error) {
	decodedPassword, err := helpers.CaesarShift(password, -svc.shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return nil, ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decodedPassword), svc.passwordSalt)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to hash user password: %v", err)
		return nil, errors.WithMessage(ErrInvalidPassword, err.Error())
	}

	return &DTO{
		Password:  string(hashedPassword),
		SecretOTP: secretOTP,
	}, nil
}

func (svc *service) ValidatePassword(ctx context.Context, credentialsDTO *DTO, password string) error {
	decodedPassword, err := helpers.CaesarShift(password, -svc.shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return ErrInvalidPassword
	}

	if err = bcrypt.CompareHashAndPassword([]byte(credentialsDTO.Password), []byte(decodedPassword)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func (svc *service) DecodePassword(ctx context.Context, password string) (string, error) {
	decodedPassword, err := helpers.CaesarShift(password, -svc.shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return "", ErrInvalidPassword
	}

	return decodedPassword, nil
}
