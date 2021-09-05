package auth

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"image/png"
	"nnw_s/config"
	"nnw_s/pkg/common"
	"nnw_s/pkg/user"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"
)

//go:generate mockgen -source=auth_service.go -destination=mocks/mock.go

type IAuthService interface {
	CreateJWTToken(email string) (string, error)
	VerifyJWTToken(id string) (*string, error)

	Generate2FaImage(email string) (*bytes.Buffer, *otp.Key, error)
	Check2FaCode(code string, secret string) bool

	CheckPassword(password string, hashPassword string) (bool, error)

	CreateTemplateUserData(secret string) (*string, error)

	CreateEmail(email string, emailType string) error
	CheckEmailCode(email string, code string, emailType string) (bool, error)
}

type AuthService struct {
	arepo AuthRepository
	urepo user.UserRepository
	cfg   config.Configurations
	emailClient config.SMTPClient
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

func NewAuthService(arepo AuthRepository, urepo user.UserRepository, cfg config.Configurations, emailClient config.SMTPClient) IAuthService {
	return &AuthService{
		arepo: arepo,
		urepo: urepo,
		cfg:   cfg,
		emailClient: emailClient,
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
	var jwtData JWTData
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

func (as *AuthService) CreateTemplateUserData(secret string) (*string, error) {
	uid := uuid.New().String()

	var templateData TemplateData
	templateData.ID = primitive.NewObjectID()
	templateData.Uid = uid
	templateData.TwoFAS = secret
	templateData.CreatedAt = time.Now()
	templateData.UpdatedAt = time.Now()

	id, err := as.arepo.CreateTemplateUserDataDb(templateData)
	if err != nil {
		return nil, err
	}

	return id, err
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

func (as *AuthService) CreateEmail(email string, emailType string) error {
	code := common.EmailCode()

	var emailData Email
	emailData.ID = primitive.NewObjectID()
	emailData.Code = code
	emailData.Email = email
	emailData.EmailType = emailType
	emailData.CreatedAt = time.Now()
	emailData.UpdatedAt = time.Now()

	err := as.arepo.CreateEmailDb(emailData)
	if err != nil {
		return err
	}

	// Set up Email template and Send email
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	t, err := template.ParseFiles(path.Join(dir, "templates/verifyTemplate.html"))
	if err != nil {
		fmt.Println(3, err)
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
		fmt.Println(1, err)
		return err
	}

	err = as.emailClient.SendMail(as.cfg.EmailFrom, []string{email}, body.Bytes())
	if err != nil {
		fmt.Println(2, err)
		return err
	}

	return nil
}

func (as *AuthService) CheckEmailCode(email string, code string, emailType string) (bool, error) {
	_, err := as.arepo.GetEmailDb(email, code, emailType)
	if err != nil {
		return false, err
	}

	return true, err
}
