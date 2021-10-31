package jwt

import (
	"context"
	"nnw_s/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const jwtExpiry = time.Second * 600

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	CreateJWT(ctx context.Context, email string) (*DTO, error)
	VerifyJWT(ctx context.Context, id string) (*Payload, error)
	DeleteJWT(ctx context.Context, token string) error
}

type service struct {
	repo      Repository
	secretKey string
}

func NewService(repo Repository, secretKey string) (Service, error) {
	if repo == nil {
		return nil, errors.NewInternal("invalid jwt repository")
	}
	if secretKey == "" {
		return nil, errors.NewInternal("invalid jwt secret key")
	}
	return &service{repo: repo, secretKey: secretKey}, nil
}

func (svc *service) CreateJWT(ctx context.Context, email string) (*DTO, error) {
	// create JWT
	payload := &Payload{
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(jwtExpiry),
	}

	// sign token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(svc.secretKey))
	if err != nil {
		return nil, ErrTokenDoesNotValid
	}

	// create JWT
	jwtData := NewJWT(signedToken)

	// save in storage
	id, err := svc.repo.SaveJWT(ctx, jwtData)
	if err != nil {
		return nil, err
	}
	return &DTO{
		ID:       id,
		Token:    signedToken,
		ExpireAt: payload.ExpiredAt,
	}, nil
}

func (svc *service) VerifyJWT(ctx context.Context, token string) (*Payload, error) {
	// get JWT from storage
	tokenDTO, err := svc.repo.GetJWT(ctx, token)
	if err != nil {
		return nil, ErrTokenDoesNotValid
	}

	// verify JWT
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenDoesNotValid
		}
		return []byte(svc.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenDTO.Jwt, &Payload{}, keyFunc)
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return nil, ErrTokenHasBeenExpired
		}
		return nil, ErrTokenDoesNotValid
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenDoesNotValid
	}
	return payload, nil
}

func (svc *service) DeleteJWT(ctx context.Context, token string) error {
	err := svc.repo.DeleteJWT(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
