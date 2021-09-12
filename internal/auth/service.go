package auth

import (
	"context"
	"nnw_s/internal/auth/mfa"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/notificator"
	"nnw_s/internal/user"
	"nnw_s/pkg/errors"
	"time"

	"github.com/sirupsen/logrus"
)

const jwtExpiry = time.Second * 60

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go

type Service interface {
	Login(ctx context.Context, dto *LoginDTO) (*TokenDTO, error)
	RegisterUser(ctx context.Context, dto *RegisterUserDTO) error
	VerifyUser(ctx context.Context, dto *VerifyUserDTO) error
	ResendVerificationEmail(ctx context.Context, email string) error
	SetupMFA(ctx context.Context, dto *SetupMfaDTO) ([]byte, error)
	ActivateUser(ctx context.Context, dto *ActivateUserDTO) error
}

type service struct {
	authRepo        Repository
	userSvc         user.Service
	notificatorSvc  notificator.Service
	verificationSvc verification.Service
	mfaSvc          mfa.Service

	log         *logrus.Logger
	emailSender string
}

type ServiceDeps struct {
	AuthRepository      Repository
	UserService         user.Service
	NotificatorService  notificator.Service
	VerificationService verification.Service
	MFAService          mfa.Service
}

func NewService(log *logrus.Logger, emailSender string, deps *ServiceDeps) (Service, error) {
	if deps == nil {
		return nil, errors.NewInternal("invalid service dependencies")
	}
	if deps.AuthRepository == nil {
		return nil, errors.NewInternal("invalid auth service")
	}
	if deps.UserService == nil {
		return nil, errors.NewInternal("invalid user service")
	}
	if deps.NotificatorService == nil {
		return nil, errors.NewInternal("invalid notificator service")
	}
	if deps.VerificationService == nil {
		return nil, errors.NewInternal("invalid verification service")
	}
	if deps.MFAService == nil {
		return nil, errors.NewInternal("invalid MFA service")
	}
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if emailSender == "" {
		return nil, errors.NewInternal("invalid sender's email")
	}
	return &service{
		authRepo:        deps.AuthRepository,
		userSvc:         deps.UserService,
		notificatorSvc:  deps.NotificatorService,
		verificationSvc: deps.VerificationService,
		log:             log,
		emailSender:     emailSender,
	}, nil
}

func (svc *service) Login(ctx context.Context, dto *LoginDTO) (*TokenDTO, error) {}

func (svc *service) RegisterUser(ctx context.Context, dto *RegisterUserDTO) error {
	// create user if not exists
	_, err := svc.userSvc.CreateUser(ctx, &user.CreateUserDTO{Email: dto.Email,
		Password: dto.Password})

	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to register user: %v", err)
		return err
	}

	// create verification code for further activation by email
	newVerificationCode, err := svc.verificationSvc.CreateVerificationCode(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create verification code: %v", err)
		return err
	}

	emailData := notificator.Email{
		Subject:   "Verify email.",
		Recipient: dto.Email,
		Sender:    svc.emailSender,
		Data: map[string]interface{}{
			"code": newVerificationCode,
		},
	}

	// send email to recipient
	if err := svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
	}

	svc.log.WithContext(ctx).Infof("verification code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *service) VerifyUser(ctx context.Context, dto *VerifyUserDTO) error {
	// check if user's verification code is valid
	if err := svc.verificationSvc.CheckVerificationCode(ctx, dto.Email, dto.Code); err != nil {
		return err
	}

	// get not activated user
	notActivatedUser, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	// activating user
	notActivatedUser.SetToVerified()

	if err = svc.userSvc.UpdateUser(ctx, notActivatedUser); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user's status field: %v", err)
		return err
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully verified", notActivatedUser.Email)
	return nil
}

func (svc *service) ResendVerificationEmail(ctx context.Context, email string) error {}

func (svc *service) SetupMFA(ctx context.Context, dto *SetupMfaDTO) ([]byte, error) {
	// find disabled user
	disabledUser, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	// generate 2FA/MFA Image
	buffImg, key, err := svc.mfaSvc.GenerateMFAImage(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create MFA image for '%s': %v", dto.Email, err)
		return nil, err
	}

	// update user SecretOTP key
	disabledUser.Credentials.SetSecretOTP(key)
	disabledUser.UpdatedAt = time.Now()

	if err = svc.userSvc.UpdateUser(ctx, disabledUser); err != nil {
		return nil, err
	}
	return buffImg.Bytes(), nil
}

func (svc *service) ActivateUser(ctx context.Context, dto *ActivateUserDTO) error {
	// find disable user
	disabledUser, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	// check 2FA/MFA Code
	if err = svc.mfaSvc.CheckMFACode(ctx, dto.Code, *disabledUser.Credentials.SecretOTP); err != nil {
		return err
	}

	// activate user
	disabledUser.SetToActive()

	if err = svc.userSvc.UpdateUser(ctx, disabledUser); err != nil {
		return err
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully activated MFA authentication", disabledUser.Email)
	return nil
}

// func (s *service) CreateJWTToken(ctx context.Context, email string) (string, error) {
// 	// Create JWT
// 	payload := &Payload{
// 		Email:     email,
// 		IssuedAt:  time.Now(),
// 		ExpiredAt: time.Now().Add(jwtExpiry),
// 	}

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	signedToken, err := jwtToken.SignedString([]byte(s.cfg.JwtSecretKey))
// 	if err != nil {
// 		return "", err
// 	}

// 	// Create JWT in DataBase
// 	var jwtData JWT
// 	jwtData.ID = primitive.NewObjectID()
// 	jwtData.Jwt = signedToken
// 	jwtData.CreatedAt = time.Now()
// 	jwtData.UpdatedAt = time.Now()

// 	id, err := s.authRepo.CreateJwt(ctx, jwtData)
// 	if err != nil {
// 		return "", err
// 	}

// 	return id, nil
// }

// func (s *service) VerifyJWTToken(ctx context.Context, id string) (string, error) {
// 	// Get Jwt from DataBase
// 	token, err := s.authRepo.GetJwt(ctx, id)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Verify JWT
// 	keyFunc := func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, errors.New("token is invalid")
// 		}
// 		return []byte(s.cfg.JwtSecretKey), nil
// 	}

// 	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
// 	if err != nil {
// 		ver, ok := err.(*jwt.ValidationError)
// 		if ok && errors.Is(ver.Inner, errors.New("token has expired")) {
// 			return "", errors.New("token has expired")
// 		}
// 		return "", errors.New("token is invalid")
// 	}

// 	payload, ok := jwtToken.Claims.(*Payload)
// 	if !ok {
// 		return "", errors.New("token is invalid")
// 	}

// 	dbUser, err := s.userRepo.GetUserByEmail(ctx, payload.Email, "active")
// 	if err != nil {
// 		return "", err
// 	}

// 	return dbUser.Email, nil
// }

// func (s *service) CheckPassword(ctx context.Context, password string, hashPassword string) error {
// 	decodePassword, err := helpers.CaesarShift(password, -s.cfg.Shift)
// 	if err != nil {
// 		return err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(decodePassword))
// 	if err != nil {
// 		return errors.New("password does not valid")
// 	}

// 	return nil
// }
