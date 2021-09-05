package service

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"nnw_s/config"
	"nnw_s/pkg/common"
	"nnw_s/pkg/user/model"
	"nnw_s/pkg/user/repository"
	"strconv"
	"time"
)

type UserService struct {
	repo repository.User
	cfg  config.Configurations
}

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewUserService(repo repository.User, cfg config.Configurations) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (as *UserService) GetUserById(userId string) (*model.User, error) {
	user, err := as.repo.GetUserByIdDb(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := as.repo.GetUserByEmailDb(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *UserService) GetTemplateUserDataById(uid string) (*model.TemplateData, error) {
	templateUserData, err := as.repo.GetTemplateUserDataByIdDb(uid)
	if err != nil {
		return nil, err
	}

	return templateUserData, nil
}

func (as *UserService) CreateUser(email string, password string, OTPKey string) (*string, error) {
	shift, err := strconv.Atoi(as.cfg.Shift)
	if err != nil {
		return nil, err
	}

	salt, err := strconv.Atoi(as.cfg.PasswordSalt)
	if err != nil {
		return nil, err
	}

	decodePassword, err := common.CaesarShift(password, -shift)
	if err != nil {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(decodePassword), salt)

	var user model.User
	user.ID = primitive.NewObjectID()
	user.Email = email
	user.Password = string(hashPassword)
	user.SecretOTPKey = OTPKey
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	id, err := as.repo.CreateUserDb(user)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (as *UserService) CreateTemplateUserData(secret string) (*string, error) {
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
