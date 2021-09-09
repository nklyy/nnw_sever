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
	CreateUser(ctx context.Context, email, password, otpSecret string) (string, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
	opts *ServiceOptions
}

type ServiceOptions struct {
	log          *logrus.Logger
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
	if opts.log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if opts.PasswordSalt == 0 {
		return nil, errors.NewInternal("invalid password salt")
	}
	if opts.Shift == 0 {
		return nil, errors.NewInternal("invalid shift")
	}
	return &service{repo: repo, opts: opts, log: opts.log}, nil
}

func (svc *service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return svc.repo.GetUserByID(ctx, userID)
}

func (svc *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return svc.repo.GetUserByEmail(ctx, email)
}

func (svc *service) CreateUser(ctx context.Context, email, password, otpSecret string) (string, error) {
	decodedPassword, err := helpers.CaesarShift(password, -svc.opts.Shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(decodedPassword), svc.opts.PasswordSalt)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to hash password: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	newUser, err := NewUser(email, string(hashPassword), otpSecret)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user entity due to validation error: %v", err)
		return "", err
	}

	id, err := svc.repo.SaveUser(ctx, newUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save user in db: %v", err)
		return "", err
	}
	return id, err
}
