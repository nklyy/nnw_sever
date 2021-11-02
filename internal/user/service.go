package user

import (
	"context"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/wallet"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserByID(ctx context.Context, userID string) (*DTO, error)
	GetUserByEmail(ctx context.Context, email string) (*DTO, error)

	CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error)

	UpdateUser(ctx context.Context, dto *DTO) error

	DeleteUserByEmail(ctx context.Context, email string) error
}

type service struct {
	repo           Repository
	credentialsSvc credentials.Service
	log            *logrus.Logger
}

func NewService(repo Repository, credentialsSvc credentials.Service, log *logrus.Logger) (Service, error) {
	if repo == nil {
		return nil, errors.NewInternal("invalid repo")
	}
	if credentialsSvc == nil {
		return nil, errors.NewInternal("invalid credentials service")
	}
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	return &service{repo: repo, credentialsSvc: credentialsSvc, log: log}, nil
}

func (svc *service) GetUserByID(ctx context.Context, userID string) (*DTO, error) {
	u, err := svc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return MapToDTO(u), nil
}

func (svc *service) GetUserByEmail(ctx context.Context, email string) (*DTO, error) {
	u, err := svc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return MapToDTO(u), nil
}

func (svc *service) CreateUser(ctx context.Context, dto *CreateUserDTO) (string, error) {
	// create user credentials
	userCredentialsDTO, err := svc.credentialsSvc.CreateCredentials(ctx, dto.Password, credentials.NilSecretOTP)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user credentials: %v", err)
		return "", err
	}

	// map credentialsDTO to entity
	userCredentials := credentials.MapToEntity(userCredentialsDTO)

	// create user with new credentials
	newUser, err := NewUser(dto.Email, wallet.NilWallet, userCredentials)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user due to validation error: %v", err)
		return "", err
	}

	id, err := svc.repo.SaveUser(ctx, newUser)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save user in db: %v", err)
		return "", err
	}
	return id, err
}

func (svc *service) UpdateUser(ctx context.Context, userDTO *DTO) error {
	// map dto to user entity
	updateUser, err := MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// update user in storage by email
	if err = svc.repo.UpdateUser(ctx, updateUser); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save user in db: %v", err)
		return err
	}
	return nil
}

func (svc *service) DeleteUserByEmail(ctx context.Context, email string) error {
	err := svc.repo.DeleteUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	return nil
}
