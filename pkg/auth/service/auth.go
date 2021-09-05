package service

import (
	"bytes"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"image/png"
	"nnw_s/config"
	"nnw_s/pkg/auth/model"
	"nnw_s/pkg/auth/repository"
	"nnw_s/pkg/common"
	repository2 "nnw_s/pkg/user/repository"
	"strconv"
	"time"
)

type AuthService struct {
	arepo repository.Authorization
	urepo repository2.User
	cfg   config.Configurations
}

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}

func NewAuthService(arepo repository.Authorization, urepo repository2.User, cfg config.Configurations) *AuthService {
	return &AuthService{
		arepo: arepo,
		urepo: urepo,
		cfg:   cfg,
	}
}

func (as *AuthService) Generate2FaImage(email string) (*bytes.Buffer, *otp.Key, error) {
	// Generate 2FA Image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NNW",
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
	err = png.Encode(&bufImage, img)
	if err != nil {
		return nil, nil, err
	}

	return &bufImage, key, nil
}

func (as *AuthService) Check2FaCode(code string, secret string) bool {
	return totp.Validate(code, secret)
}

func (as *AuthService) CreateJWTToken(email string) (string, error) {
	// Create JWT
	payload := &Payload{
		Email:     email,
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

	id, err := as.arepo.CreateJwtDb(jwtData)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (as *AuthService) VerifyJWTToken(id string) (*string, error) {
	// Get Jwt from DataBase
	token, err := as.arepo.GetJwtDb(id)
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

	user, err := as.urepo.GetUserByEmailDb(payload.Email)
	if err != nil {
		return nil, err
	}

	return &user.Email, nil
}

func (as *AuthService) CheckPassword(password string, hashPassword string) (bool, error) {
	shift, err := strconv.Atoi(as.cfg.Shift)
	if err != nil {
		return false, err
	}

	decodePassword, err := common.CaesarShift(password, -shift)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(decodePassword))
	if err != nil {
		return false, nil
	}

	return true, nil
}
