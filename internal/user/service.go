package user

import (
	"context"
	"nnw_s/config"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string, status string) (*User, error)

	CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error)

	UpdateUser(ctx context.Context, dto *User) error
	UpdateDisableUser(ctx context.Context, oldUser *User) error
}

type service struct {
	repo Repository
	log  *logrus.Logger
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config, log *logrus.Logger) (Service, error) {
	if repo == nil {
		return nil, errors.NewInternal("invalid repo")
	}

	if cfg == nil {
		return nil, errors.NewInternal("invalid service options")
	}

	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}

	if cfg.PasswordSalt == 0 {
		return nil, errors.NewInternal("invalid password salt")
	}
	if cfg.Shift == 0 {
		return nil, errors.NewInternal("invalid shift")
	}

	return &service{repo: repo, cfg: cfg, log: log}, nil
}

func (svc *service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return svc.repo.GetUserByID(ctx, userID)
}

func (svc *service) GetUserByEmail(ctx context.Context, email string, status string) (*User, error) {
	return svc.repo.GetUserByEmail(ctx, email, status)
}

func (svc *service) CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error) {
	decodedPassword, err := helpers.CaesarShift(dto.Password, -svc.cfg.Shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	newUser, err := NewUser(dto.Email, decodedPassword, dto.SecretOTP)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user due to validation error: %v", err)
		return "", err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), svc.cfg.PasswordSalt)
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

func (svc *service) UpdateUser(ctx context.Context, updateUser *User) error {
	err := svc.repo.UpdateUser(ctx, updateUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save user in db: %v", err)
		return err
	}
	return nil
}

func (svc *service) UpdateDisableUser(ctx context.Context, oldUser *User) error {
	decodedPassword, err := helpers.CaesarShift(oldUser.Password, -svc.cfg.Shift)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to decode password: %v", err)
		return errors.NewInternal(err.Error())
	}

	newDisableUser, err := NewDisableUser(oldUser.Email, decodedPassword, oldUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user due to validation error: %v", err)
		return err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newDisableUser.Password), svc.cfg.PasswordSalt)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to hash user password: %v", err)
		return errors.NewInternal(err.Error())
	}

	newDisableUser.Password = string(hashPassword)

	err = svc.repo.UpdateDisableUser(ctx, newDisableUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user in db: %v", err)
		return err
	}
	return nil
}
