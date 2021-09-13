package jwt

import (
	"context"
	"nnw_s/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const jwtExpiry = time.Second * 60

type Service interface {
	CreateJWT(ctx context.Context, email string) (*DTO, error)
	VerifyJWT(ctx context.Context, id string) error
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

func (svc *service) VerifyJWT(ctx context.Context, id string) error {
	// get JWT from storage
	token, err := svc.repo.GetJWT(ctx, id)
	if err != nil {
		return ErrTokenDoesNotValid
	}

	// verify JWT
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenDoesNotValid
		}
		return []byte(svc.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token.Jwt, &Payload{}, keyFunc)
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return ErrTokenHasBeenExpired
		}
		return ErrTokenDoesNotValid
	}

	_, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return ErrTokenDoesNotValid
	}
	return nil
}