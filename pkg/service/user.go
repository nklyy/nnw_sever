package service

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
	"time"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetUserById(userId string) (*model.User, error) {
	user, err := us.repo.GetUserByIdDb(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByLogin(login string) (*model.User, error) {
	user, err := us.repo.GetUserByLoginDb(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetTemplateUserDataById(uid string) (*model.TemplateData, error) {
	templateUserData, err := us.repo.GetTemplateUserDataByIdDb(uid)
	if err != nil {
		return nil, err
	}

	return templateUserData, nil
}

func (us *UserService) CreateUser(login string, email string, password string, OTPKey string) (*string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	var user model.User
	user.ID = primitive.NewObjectID()
	user.Login = login
	user.Email = email
	user.Password = string(hashPassword)
	user.SecretOTPKey = OTPKey

	id, err := us.repo.CreateUserDb(user)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (us *UserService) CreateTemplateUserData(secret string) (*string, error) {
	uid := uuid.New().String()

	var templateData model.TemplateData
	templateData.ID = primitive.NewObjectID()
	templateData.Uid = uid
	templateData.TwoFAS = secret
	templateData.CreatedAt = time.Now()
	templateData.UpdatedAt = time.Now()

	id, err := us.repo.CreateTemplateUserDataDb(templateData)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (us *UserService) UpdateUser(user model.User) error {
	err := us.repo.UpdateUserDb(user)
	if err != nil {
		return err
	}

	return nil
}
