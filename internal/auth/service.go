package auth

import (
	"context"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/mfa"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/notificator"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"time"

	"github.com/sirupsen/logrus"
)

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
	userSvc         user.Service
	notificatorSvc  notificator.Service
	verificationSvc verification.Service
	mfaSvc          mfa.Service
	jwtSvc          jwt.Service
	credentialsSvc  credentials.Service

	log         *logrus.Logger
	emailSender string
}

type ServiceDeps struct {
	UserService         user.Service
	NotificatorService  notificator.Service
	VerificationService verification.Service
	MFAService          mfa.Service
	JWTService          jwt.Service
	CredentialsService  credentials.Service
}

func NewService(log *logrus.Logger, emailSender string, deps *ServiceDeps) (Service, error) {
	if deps == nil {
		return nil, errors.NewInternal("invalid service dependencies")
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
	if deps.JWTService == nil {
		return nil, errors.NewInternal("invalid JWT service")
	}
	if deps.CredentialsService == nil {
		return nil, errors.NewInternal("invalid credentials service")
	}
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if emailSender == "" {
		return nil, errors.NewInternal("invalid sender's email")
	}
	return &service{
		userSvc:         deps.UserService,
		notificatorSvc:  deps.NotificatorService,
		verificationSvc: deps.VerificationService,
		log:             log,
		emailSender:     emailSender,
	}, nil
}

func (svc *service) Login(ctx context.Context, dto *LoginDTO) (*TokenDTO, error) {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	registeredUser, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !registeredUser.IsActive() || !registeredUser.IsVerified {
		return nil, ErrPermissionDenied
	}

	// map from entity to credentials dto
	credentialsDTO := credentials.MapToDTO(registeredUser.Credentials)

	// check password
	if err = svc.credentialsSvc.ValidatePassword(ctx, credentialsDTO, dto.Password); err != nil {
		return nil, err
	}

	// check 2FA/MFA Code
	if err = svc.mfaSvc.CheckMFACode(ctx, dto.Code, *registeredUser.Credentials.SecretOTP); err != nil {
		return nil, err
	}

	// create JWT
	jwtTokenDTO, err := svc.jwtSvc.CreateJWT(ctx, dto.Email)
	if err != nil {
		return nil, errors.WithMessage(ErrUnauthorized, err.Error())
	}
	return &TokenDTO{Token: jwtTokenDTO.Token, ExpireAt: jwtTokenDTO.ExpireAt}, nil
}

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
	if err = svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
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
	notActivatedUserDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	// mapping userDTO to user entity
	notActivatedUser, err := user.MapToEntity(notActivatedUserDTO)
	if err != nil {
		return err
	}

	// activating user
	notActivatedUser.SetToVerified()

	// map back to dto
	notActivatedUserDTO = user.MapToDTO(notActivatedUser)

	// update user entity in storage
	if err = svc.userSvc.UpdateUser(ctx, notActivatedUserDTO); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user's status field: %v", err)
		return err
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully verified", notActivatedUser.Email)
	return nil
}

// todo
func (svc *service) ResendVerificationEmail(ctx context.Context, email string) error {
	return errors.NewInternal("not implemented yet")
}

func (svc *service) SetupMFA(ctx context.Context, dto *SetupMfaDTO) ([]byte, error) {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	// map userDTO to user
	disabledUser, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, err
	}

	// check if user is active or not verified, if yes - return ErrPermissionDenied
	if disabledUser.IsActive() || !disabledUser.IsVerified {
		return nil, ErrPermissionDenied
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

	// map to userDTO
	userDTO = user.MapToDTO(disabledUser)

	// update user in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		return nil, err
	}
	return buffImg.Bytes(), nil
}

func (svc *service) ActivateUser(ctx context.Context, dto *ActivateUserDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	// map userDTO to user
	disabledUser, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// check if user is active or not verified, if yes - return ErrPermissionDenied
	if disabledUser.IsActive() || !disabledUser.IsVerified {
		return ErrPermissionDenied
	}

	// check 2FA/MFA Code
	if err = svc.mfaSvc.CheckMFACode(ctx, dto.Code, *disabledUser.Credentials.SecretOTP); err != nil {
		return err
	}

	// activate user
	disabledUser.SetToActive()

	// map back to DTO
	userDTO = user.MapToDTO(disabledUser)

	// update user in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		return err
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully activated MFA authentication", disabledUser.Email)
	return nil
}
