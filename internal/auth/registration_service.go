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

const (
	emailVerificationSubject      = "Verification email."
	emailVerificationTopic        = "Verification email."
	emailVerificationMessage      = "You're receiving this e-mail because you requested a verify your email for your NoName Wallet account."
	emailVerificationTemplateName = "authTemplate.html"
)

func NewRegistrationService(log *logrus.Logger, emailSender string, deps *ServiceDeps) (RegistrationService, error) {
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
	userDTO, _ := svc.userSvc.GetUserByEmail(ctx, dto.Email)

	if userDTO == nil {
		_, err := svc.userSvc.CreateUser(ctx, &user.CreateUserDTO{Email: dto.Email, Password: dto.Password})
		if err != nil {
			svc.log.WithContext(ctx).Errorf("failed to register user: %v", err)
			return err
		}
	} else if userDTO.Status == "disabled" {
		err := svc.userSvc.DeleteUserByEmail(ctx, dto.Email)
		if err != nil {
			svc.log.WithContext(ctx).Errorf("failed to delete user: %v", err)
			return err
		}

		_, err = svc.userSvc.CreateUser(ctx, &user.CreateUserDTO{Email: dto.Email, Password: dto.Password})
		if err != nil {
			svc.log.WithContext(ctx).Errorf("failed to register user: %v", err)
			return err
		}
	} else {
		svc.log.WithContext(ctx).Errorf("failed to register user %v: user already exists", dto.Email)
		return user.ErrAlreadyExists
	}

	// create verification code for further activation by email
	newVerificationCode, err := svc.verificationSvc.CreateVerificationCode(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create verification code: %v", err)
		return ErrFailedCreateCode
	}

	emailData := notificator.Email{
		Subject:   emailVerificationSubject,
		Recipient: dto.Email,
		Sender:    svc.emailSender,
		Template:  emailVerificationTemplateName,
		Data: map[string]interface{}{
			"topic":   emailVerificationTopic,
			"message": emailVerificationMessage,
			"code":    newVerificationCode,
		},
	}

	// send email to recipient
	if err = svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
		return ErrFailedSendEmail
	}

	svc.log.WithContext(ctx).Infof("verification code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *registrationSvc) VerifyUser(ctx context.Context, dto *VerifyUserDTO) error {
	// check if user's verification code is valid
	if err := svc.verificationSvc.CheckVerificationCode(ctx, dto.Email, dto.Code); err != nil {
		return ErrInvalidCode
	}

	// get not activated user
	notActivatedUserDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	if notActivatedUserDTO.IsVerified == true {
		svc.log.WithContext(ctx).Errorf("user %v, already verify", dto.Email)
		return user.ErrUserAlreadyVerify
	}

	// mapping userDTO to user entity
	notActivatedUser, err := user.MapToEntity(notActivatedUserDTO)
	if err != nil {
		return ErrInvalidDTO
	}

	// activating user
	notActivatedUser.SetToVerified()

	// map back to dto
	notActivatedUserDTO = user.MapToDTO(notActivatedUser)

	// update user entity in storage
	if err = svc.userSvc.UpdateUser(ctx, notActivatedUserDTO); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user's status field: %v", err)
		return user.ErrFailedUpdateUser
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully verified", notActivatedUser.Email)
	return nil
}

func (svc *registrationSvc) ResendVerificationEmail(ctx context.Context, dto *ResendActivationEmailDTO) error {
	// get not activated user
	notActivatedUserDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// mapping userDTO to user entity
	userEntity, err := user.MapToEntity(notActivatedUserDTO)
	if err != nil {
		return ErrInvalidDTO
	}

	// if user does not active or not verified return ErrPermissionDenied
	if userEntity.IsActive() || userEntity.IsVerified {
		return user.ErrUserAlreadyVerify
	}

	// create verification code for further activation by email
	newVerificationCode, err := svc.verificationSvc.CreateVerificationCode(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create verification code: %v", err)
		return ErrFailedCreateCode
	}

	emailData := notificator.Email{
		Subject:   emailVerificationSubject,
		Recipient: dto.Email,
		Sender:    svc.emailSender,
		Template:  emailVerificationTemplateName,
		Data: map[string]interface{}{
			"topic":   emailVerificationTopic,
			"message": emailVerificationMessage,
			"code":    newVerificationCode,
		},
	}

	// send email to recipient
	if err = svc.notificatorSvc.SendEmail(ctx, &emailData); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
		return ErrFailedSendEmail
	}

	svc.log.WithContext(ctx).Infof("verification code successfully sent to: %s", dto.Email)
	return nil
}

func (svc *registrationSvc) SetupTwoFA(ctx context.Context, dto *SetupTwoFaDTO) ([]byte, error) {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, user.ErrNotFound
	}

	// map userDTO to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, ErrInvalidDTO
	}

	// check if user is active or not verified, if yes - return ErrPermissionDenied
	if userEntity.IsActive() || !userEntity.IsVerified {
		return nil, user.ErrUserAlreadyVerify
	}

	// generate TwoFa Image
	buffImg, key, err := svc.twoFaSvc.GenerateTwoFAImage(ctx, dto.Email)
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to create TwoFA image for '%s': %v", dto.Email, err)
		return nil, ErrFailedGenerateTwoFaImage
	}

	// update user SecretOTP key
	userEntity.Credentials.SetSecretOTP(key)
	userEntity.UpdatedAt = time.Now()

	// map to userDTO
	userDTO = user.MapToDTO(userEntity)

	// update user in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		return nil, user.ErrFailedUpdateUser
	}
	return buffImg.Bytes(), nil
}

func (svc *registrationSvc) ActivateUser(ctx context.Context, dto *ActivateUserDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return user.ErrNotFound
	}

	// map userDTO to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return ErrInvalidDTO
	}

	// check if user is active or not verified, if yes - return ErrPermissionDenied
	if userEntity.IsActive() && userEntity.IsVerified {
		return user.ErrUserAlreadyActive
	}

	// check TwoFA Code
	if err = svc.twoFaSvc.CheckTwoFACode(ctx, dto.Code, *userEntity.Credentials.SecretOTP); err != nil {
		return ErrInvalidCode
	}

	// activate user
	userEntity.SetToActive()

	// map back to DTO
	userDTO = user.MapToDTO(userEntity)

	// update user in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		return user.ErrFailedUpdateUser
	}

	svc.log.WithContext(ctx).Infof("user '%s' successfully activated TwoFA authentication", userEntity.Email)
	return nil
}
