package auth

import (
	"context"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/user"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/notificator"
	"time"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=registration_service.go -destination=mocks/registration_service_mock.go
type RegistrationService interface {
	RegisterUser(ctx context.Context, dto *RegisterUserDTO) error
	VerifyUser(ctx context.Context, dto *VerifyUserDTO) error
	ResendVerificationEmail(ctx context.Context, dto *ResendActivationEmailDTO) error
	SetupTwoFA(ctx context.Context, dto *SetupTwoFaDTO) ([]byte, error)
	ActivateUser(ctx context.Context, dto *ActivateUserDTO) error
}

type registrationSvc struct {
	userSvc         user.Service
	notificatorSvc  notificator.Service
	verificationSvc verification.Service
	twoFaSvc        twofa.Service

	log         *logrus.Logger
	emailSender string
}

func NewRegistrationService(log *logrus.Logger, emailSender string, deps *ServiceDeps) (RegistrationService, error) {
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
	if deps.TwoFAService == nil {
		return nil, errors.NewInternal("invalid TwoFA service")
	}
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if emailSender == "" {
		return nil, errors.NewInternal("invalid sender's email")
	}
	return &registrationSvc{
		userSvc:         deps.UserService,
		notificatorSvc:  deps.NotificatorService,
		verificationSvc: deps.VerificationService,
		log:             log,
		emailSender:     emailSender,
		twoFaSvc:        deps.TwoFAService,
	}, nil
}

func (svc *registrationSvc) RegisterUser(ctx context.Context, dto *RegisterUserDTO) error {
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
		return err
	}

	svc.log.WithContext(ctx).Infof("verification code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *registrationSvc) VerifyUser(ctx context.Context, dto *VerifyUserDTO) error {
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

func (svc *registrationSvc) ResendVerificationEmail(ctx context.Context, dto *ResendActivationEmailDTO) error {
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

	// check if user not active and not verified
	if !notActivatedUser.IsActive() && !notActivatedUser.IsVerified {
		return ErrPermissionDenied
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
		return err
	}

	svc.log.WithContext(ctx).Infof("verification code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *registrationSvc) SetupTwoFA(ctx context.Context, dto *SetupTwoFaDTO) ([]byte, error) {
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

	// generate TwoFa Image
	buffImg, key, err := svc.twoFaSvc.GenerateTwoFAImage(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create TwoFA image for '%s': %v", dto.Email, err)
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

func (svc *registrationSvc) ActivateUser(ctx context.Context, dto *ActivateUserDTO) error {
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

	// check TwoFA Code
	if err = svc.twoFaSvc.CheckTwoFACode(ctx, dto.Code, *disabledUser.Credentials.SecretOTP); err != nil {
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

	svc.log.WithContext(ctx).Infof("user '%s' successfully activated TwoFA authentication", disabledUser.Email)
	return nil
}
