package user

import (
	"context"

	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
	opts *ServiceOptions
}

type ServiceOptions struct {
	Log          *logrus.Logger
	Shift        int
	PasswordSalt int
}

func NewService(repo Repository, opts *ServiceOptions) (Service, error) {
	if repo == nil {
		return nil, errors.NewInternal("invalid repo")
	}
	if opts == nil {
		return nil, errors.NewInternal("invalid service options")
	}
	if opts.Log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if opts.PasswordSalt == 0 {
		return nil, errors.NewInternal("invalid password salt")
	}
	if opts.Shift == 0 {
		return nil, errors.NewInternal("invalid shift")
	}
	return &service{repo: repo, opts: opts, log: opts.Log}, nil
}

func (svc *service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return svc.repo.GetUserByID(ctx, userID)
}

func (svc *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return svc.repo.GetUserByEmail(ctx, email)
}

func (svc *service) CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error) {
	decodedPassword, err := helpers.CaesarShift(dto.Password, -svc.opts.Shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	newUser, err := NewUser(dto.Email, decodedPassword, dto.SecretOTP)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user due to validation error: %v", err)
		return "", err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), svc.opts.PasswordSalt)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to hash user password: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	newUser.Password = string(hashPassword)

	id, err := svc.repo.SaveUser(ctx, newUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save user in db: %v", err)
		return "", err
	}
	return id, err
}
