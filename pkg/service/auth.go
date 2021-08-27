package service

import (
	"bytes"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"image/png"
	"nnw_s/config"
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
	"strconv"
	"time"
)

type AuthService struct {
	repo repository.Authorization
	cfg  config.Configurations
}

type Payload struct {
	Login     string    `json:"login"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}

func NewAuthService(repo repository.Authorization, cfg config.Configurations) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
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
	shift, err := strconv.Atoi(as.cfg.Shift)
	if err != nil {
		return nil, err
	}

	decodePassword, err := caesarShift(password, -shift)
	if err != nil {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(decodePassword), 15)

	var user model.User
	user.ID = primitive.NewObjectID()
	user.Login = login
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

func (as *AuthService) Generate2FaImage(login string) (*bytes.Buffer, *otp.Key, error) {
	// Generate 2FA Image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NNW",
		AccountName: login,
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

func (as *AuthService) CreateJWTToken(login string) (string, error) {
	// Create JWT
	payload := &Payload{
		Login:     login,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Second * 60),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(as.cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}

	// Create JWT in DataBase
	var jwtData model.JWTData
	jwtData.ID = primitive.NewObjectID()
	jwtData.Jwt = signedToken
	jwtData.CreatedAt = time.Now()
	jwtData.UpdatedAt = time.Now()

	id, err := as.repo.CreateJwtDb(jwtData)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (as *AuthService) VerifyJWTToken(id string) (*string, error) {
	// Get Jwt from DataBase
	token, err := as.repo.GetJwtDb(id)
	if err != nil {
		return nil, err
	}

	// Verify JWT
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(as.cfg.JwtSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(*token, &Payload{}, keyFunc)
	if err != nil {
		ver, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(ver.Inner, errors.New("token has expired")) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("token is invalid")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	user, err := as.repo.GetUserByLoginDb(payload.Login)
	if err != nil {
		return nil, err
	}

	return &user.Login, nil
}

func (as *AuthService) CheckPassword(password string, hashPassword string) (bool, error) {
	shift, err := strconv.Atoi(as.cfg.Shift)
	if err != nil {
		return false, err
	}

	decodePassword, err := caesarShift(password, -shift)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(decodePassword))
	if err != nil {
		return false, nil
	}

	return true, nil
}
