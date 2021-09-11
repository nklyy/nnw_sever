package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"image/png"
	"nnw_s/config"
	"nnw_s/internal/user"
	"nnw_s/pkg/helpers"
	"nnw_s/pkg/smtp"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const jwtExpiry = time.Second * 60

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go

type Service interface {
	CreateJWTToken(ctx context.Context, email string) (string, error)
	VerifyJWTToken(ctx context.Context, id string) (string, error)

	Generate2FaImage(ctx context.Context, email string) (*bytes.Buffer, *otp.Key, error)
	Check2FaCode(code string, secret string) bool

	CheckPassword(ctx context.Context, password string, hashPassword string) error

	CreateEmail(ctx context.Context, email string, emailType string) error
	CheckEmailCode(ctx context.Context, email string, code string, emailType string) (bool, error)
}

type service struct {
	authRepo    Repository
	userRepo    user.Repository
	emailClient smtp.Client
	cfg         config.Config
}

func NewService(authRepo Repository, userRepo user.Repository, cfg config.Config, emailClient smtp.Client) Service {
	return &service{
		authRepo:    authRepo,
		userRepo:    userRepo,
		cfg:         cfg,
		emailClient: emailClient,
	}
}

func (s *service) Generate2FaImage(ctx context.Context, email string) (*bytes.Buffer, *otp.Key, error) {
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

func (s *service) Check2FaCode(code, secret string) bool {
	return totp.Validate(code, secret)
}

func (s *service) CreateJWTToken(ctx context.Context, email string) (string, error) {
	// Create JWT
	payload := &Payload{
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(jwtExpiry),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(s.cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}

	// Create JWT in DataBase
	var jwtData JWT
	jwtData.ID = primitive.NewObjectID()
	jwtData.Jwt = signedToken
	jwtData.CreatedAt = time.Now()
	jwtData.UpdatedAt = time.Now()

	id, err := s.authRepo.CreateJwt(ctx, jwtData)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) VerifyJWTToken(ctx context.Context, id string) (string, error) {
	// Get Jwt from DataBase
	token, err := s.authRepo.GetJwt(ctx, id)
	if err != nil {
		return "", err
	}

	// Verify JWT
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(s.cfg.JwtSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		ver, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(ver.Inner, errors.New("token has expired")) {
			return "", errors.New("token has expired")
		}
		return "", errors.New("token is invalid")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return "", errors.New("token is invalid")
	}

	dbUser, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", err
	}

	return dbUser.Email, nil
}

func (s *service) CheckPassword(ctx context.Context, password string, hashPassword string) error {
	shift, err := strconv.Atoi(s.cfg.Shift)
	if err != nil {
		return err
	}

	decodePassword, err := helpers.CaesarShift(password, -shift)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(decodePassword))
	if err != nil {
		return errors.New("password does not valid")
	}

	return nil
}

func (s *service) CreateEmail(ctx context.Context, email, emailType string) error {
	code := helpers.EmailCode()

	var emailData Email
	emailData.ID = primitive.NewObjectID()
	emailData.Code = code
	emailData.Email = email
	emailData.EmailType = emailType
	emailData.CreatedAt = time.Now()
	emailData.UpdatedAt = time.Now()

	err := s.authRepo.CreateEmail(ctx, emailData)
	if err != nil {
		return err
	}

	// Set up Email template and Send email
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	root := filepath.Dir(dir)

	t, err := template.ParseFiles(path.Join(root, "templates/verifyTemplate.html"))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verify Email. \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Code string
	}{
		Code: code,
	})
	if err != nil {
		return err
	}

	err = s.emailClient.SendMail(s.cfg.EmailFrom, []string{email}, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CheckEmailCode(ctx context.Context, email, code, emailType string) (bool, error) {
	_, err := s.authRepo.GetEmail(ctx, email, code, emailType)
	if err != nil {
		return false, err
	}

	return true, err
}
