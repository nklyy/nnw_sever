package mfa

import (
	"bytes"
	"context"
	"image/png"
	"nnw_s/pkg/errors"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GenerateMFAImage(ctx context.Context, email string) (*bytes.Buffer, *otp.Key, error)
	CheckMFACode(ctx context.Context, code, secret string) error
}

type service struct {
	issuer string
}

func NewService(issuer string) (Service, error) {
	if issuer == "" {
		return nil, errors.NewInternal("invalid issuer")
	}
	return &service{issuer: issuer}, nil
}

func (svc *service) GenerateMFAImage(ctx context.Context, email string) (*bytes.Buffer, *otp.Key, error) {
	// Generate 2FA Image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      svc.issuer,
		AccountName: email,
	})

	if err != nil {
		return nil, nil, err
	}

	var bufImage bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, nil, err
	}

	// Encode image
	if err = png.Encode(&bufImage, img); err != nil {
		return nil, nil, err
	}
	return &bufImage, key, nil
}

func (svc *service) CheckMFACode(ctx context.Context, code, secret string) error {
	if !totp.Validate(code, secret) {
		return ErrInvalidMFACode
	}
	return nil
}
