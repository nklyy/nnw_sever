package user

import (
	"context"
	"nnw_s/config"
	"strconv"
	"time"

	"nnw_s/pkg/helpers"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//todo
// move templateUserData, jwtData, 2fa

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetTemplateUserDataByID(ctx context.Context, uid string) (*TemplateData, error)

	CreateUser(ctx context.Context, email, password, otpKey string) (string, error)
	CreateTemplateUserData(ctx context.Context, secret string) (string, error)
}

type service struct {
	repo Repository
	cfg  config.Config
}

func NewService(repo Repository, cfg config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (svc *service) GetUserByID(ctx context.Context, userId string) (*User, error) {
	user, err := svc.repo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *service) GetTemplateUserDataByID(ctx context.Context, uid string) (*TemplateData, error) {
	templateUserData, err := as.repo.GetTemplateUserDataByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return templateUserData, nil
}

func (as *service) CreateUser(ctx context.Context, email string, password string, OTPKey string) (string, error) {
	shift, err := strconv.Atoi(as.cfg.Shift)
	if err != nil {
		return "", err
	}

	salt, err := strconv.Atoi(as.cfg.PasswordSalt)
	if err != nil {
		return "", err
	}

	decodePassword, err := helpers.CaesarShift(password, -shift)
	if err != nil {
		return "", err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(decodePassword), salt)

	var user User
	user.ID = primitive.NewObjectID()
	user.Email = email
	user.Password = string(hashPassword)
	user.SecretOTPKey = OTPKey
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	id, err := as.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, err
}

func (as *service) CreateTemplateUserData(ctx context.Context, secret string) (string, error) {
	uid := uuid.New().String()

	var templateData TemplateData
	templateData.ID = primitive.NewObjectID()
	templateData.Uid = uid
	templateData.TwoFAS = secret
	templateData.CreatedAt = time.Now()
	templateData.UpdatedAt = time.Now()

	id, err := as.repo.CreateTemplateUserData(ctx, templateData)
	if err != nil {
		return "", err
	}

	return id, err
}
