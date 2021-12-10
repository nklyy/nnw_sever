package auth

import (
	"context"
	"github.com/sirupsen/logrus"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/notificator"
)

//go:generate mockgen -source=resetPassword_service.go -destination=mocks/resetPassword_service_mock.go
type ResetPasswordService interface {
	ResetPassword(ctx context.Context, dto *ResetPasswordDTO) error
	ResendResetPasswordEmail(ctx context.Context, dto *ResendResetPasswordDTO) error
	ResetPasswordCode(ctx context.Context, dto *ResetPasswordCodedDTO) error
	SetupNewPassword(ctx context.Context, dto *SetupNewPasswordDTO) error
}

type resetPasswordSvc struct {
	userSvc         user.Service
	notificatorSvc  notificator.Service
	verificationSvc verification.Service
	credentialsSvc  credentials.Service

	log         *logrus.Logger
	emailSender string
}

const (
	emailResetPasswordSubject      = "Reset password."
	emailResetPasswordTopic        = "Reset password."
	emailResetPasswordMessage      = "You're receiving this e-mail because you requested a reset your password for your NoName Wallet account."
	emailResetPasswordTemplateName = "authTemplate.html"
)

func NewResetPasswordService(log *logrus.Logger, emailSender string, deps *ServiceDeps) (ResetPasswordService, error) {
	if deps == nil {
		return nil, errors.NewInternal("invalid service dependencies")
	}
	if deps.UserService == nil {
		return nil, errors.NewInternal("invalid user service")
	}
	if deps.NotificatorService == nil {
		return nil, errors.NewInternal("invalid notification service")
	}
	if deps.VerificationService == nil {
		return nil, errors.NewInternal("invalid verification service")
	}
	if deps.TwoFAService == nil {
		return nil, errors.NewInternal("invalid TwoFA service")
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

	return &resetPasswordSvc{
		userSvc:         deps.UserService,
		notificatorSvc:  deps.NotificatorService,
		verificationSvc: deps.VerificationService,
		credentialsSvc:  deps.CredentialsService,
		log:             log,
		emailSender:     emailSender,
	}, nil
}

func (svc *resetPasswordSvc) ResetPassword(ctx context.Context, dto *ResetPasswordDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !userEntity.IsActive() || !userEntity.IsVerified {
		return ErrPermissionDenied
	}

	newResetPasswordCode, err := svc.verificationSvc.CreateResetPasswordCode(ctx, userEntity.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create reset password code: %v", err)
		return err
	}

	emailData := notificator.Email{
		Subject:   emailResetPasswordSubject,
		Recipient: dto.Email,
		Sender:    svc.emailSender,
		Template:  emailResetPasswordTemplateName,
		Data: map[string]interface{}{
			"topic":   emailResetPasswordTopic,
			"message": emailResetPasswordMessage,
			"code":    newResetPasswordCode,
		},
	}

	// send email to recipient
	if err = svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
		return err
	}

	svc.log.WithContext(ctx).Infof("reset password code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *resetPasswordSvc) ResendResetPasswordEmail(ctx context.Context, dto *ResendResetPasswordDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !userEntity.IsActive() || !userEntity.IsVerified {
		return ErrPermissionDenied
	}

	newResetPasswordCode, err := svc.verificationSvc.CreateResetPasswordCode(ctx, userEntity.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create reset password code: %v", err)
		return err
	}

	emailData := notificator.Email{
		Subject:   emailResetPasswordSubject,
		Recipient: dto.Email,
		Sender:    svc.emailSender,
		Template:  emailResetPasswordTemplateName,
		Data: map[string]interface{}{
			"topic":   emailResetPasswordTopic,
			"message": emailResetPasswordMessage,
			"code":    newResetPasswordCode,
		},
	}

	// send email to recipient
	if err = svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
		return err
	}

	svc.log.WithContext(ctx).Infof("reset password code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *resetPasswordSvc) ResetPasswordCode(ctx context.Context, dto *ResetPasswordCodedDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !userEntity.IsActive() || !userEntity.IsVerified {
		return ErrPermissionDenied
	}

	err = svc.verificationSvc.CheckResetPasswordCode(ctx, dto.Email, dto.Code)
	if err != nil {
		return err
	}

	return nil
}

func (svc *resetPasswordSvc) SetupNewPassword(ctx context.Context, dto *SetupNewPasswordDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !userEntity.IsActive() || !userEntity.IsVerified {
		return ErrPermissionDenied
	}

	// Create new credentials
	userCredentialsDTO, err := svc.credentialsSvc.CreateCredentials(ctx, dto.Password, userEntity.Credentials.SecretOTP)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create user credentials: %v", err)
		return err
	}

	// map credentialsDTO to entity
	userCredentials := credentials.MapToEntity(userCredentialsDTO)

	// Set-up new credentials
	userEntity.Credentials = userCredentials

	// map to userDTO
	resetPasswordUserDTO := user.MapToDTO(userEntity)

	err = svc.userSvc.UpdateUser(ctx, resetPasswordUserDTO)
	if err != nil {
		return err
	}

	return nil
}
