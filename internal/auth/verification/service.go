package verification

import (
	"context"
	"nnw_s/pkg/errors"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	CheckVerificationCode(ctx context.Context, email, code string) error
	CreateVerificationCode(ctx context.Context, email string) (string, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(repo Repository, log *logrus.Logger) (Service, error) {
	if repo == nil {
		return nil, errors.NewInternal("invalid repo")
	}
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	return &service{repo: repo, log: log}, nil
}

func (svc *service) CreateVerificationCode(ctx context.Context, email string) (string, error) {
	newCode, err := NewCode(email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create verification code: %v", err)
		return "", err
	}
	if err := svc.repo.SaveVerificationCode(ctx, newCode); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to save verification code in db: %v", err)
		return "", err
	}
	return newCode.Code, nil
}

func (svc *service) CheckVerificationCode(ctx context.Context, email, code string) error {
	_, err := svc.repo.GetVerificationCode(ctx, email, code)
	if err != nil {
		if err == ErrCodeNotFound {
			return ErrInvalidCode
		}
		return err
	}
	return nil
}
