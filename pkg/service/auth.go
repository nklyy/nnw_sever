package service

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"image/png"
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) GetUserById(userId string) (*model.User, error) {
	user, err := as.repo.GetUserByIdDb(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *AuthService) GetUserByLogin(login string) (*model.User, error) {
	user, err := as.repo.GetUserByLoginDb(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *AuthService) GetTemplateUserDataById(uid string) (*model.TemplateData, error) {
	templateUserData, err := as.repo.GetTemplateUserDataByIdDb(uid)
	if err != nil {
		return nil, err
	}

	return templateUserData, nil
}

func (as *AuthService) CreateUser(login string, email string, password string, OTPKey string) (*string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	var user model.User
	user.ID = primitive.NewObjectID()
	user.Login = login
	user.Email = email
	user.Password = string(hashPassword)
	user.SecretOTPKey = OTPKey

	id, err := as.repo.CreateUserDb(user)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (as *AuthService) CreateTemplateUserData(secret string) (*string, error) {
	uid := uuid.New().String()

	var templateData model.TemplateData
	templateData.ID = primitive.NewObjectID()
	templateData.Uid = uid
	templateData.TwoFAS = secret
	templateData.CreatedAt = time.Now()
	templateData.UpdatedAt = time.Now()

	id, err := as.repo.CreateTemplateUserDataDb(templateData)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (as *AuthService) Generate2FaImage() (*bytes.Buffer, *otp.Key, error) {
	// Generate 2FA Image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NNW",
		AccountName: "example@examole.com",
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
	err = png.Encode(&bufImage, img)
	if err != nil {
		return nil, nil, err
	}

	return &bufImage, key, nil
}
