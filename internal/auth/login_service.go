package auth

import (
	"context"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/notificator"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=login_service.go -destination=mocks/login_service_mock.go
type LoginService interface {
	Login(ctx context.Context, dto *LoginDTO) error
	CheckCode(ctx context.Context, dto *LoginCodeDTO) (*TokenDTO, error)
	Logout(ctx context.Context, email string) error
}

type loginSvc struct {
	userSvc        user.Service
	twoFaSvc       twofa.Service
	jwtSvc         jwt.Service
	credentialsSvc credentials.Service

	log *logrus.Logger
}

type ServiceDeps struct {
	UserService         user.Service
	NotificatorService  notificator.Service
	VerificationService verification.Service
	TwoFAService        twofa.Service
	JWTService          jwt.Service
	CredentialsService  credentials.Service
}

func NewLoginService(log *logrus.Logger, deps *ServiceDeps) (LoginService, error) {
	if deps == nil {
		return nil, errors.NewInternal("invalid service dependencies")
	}
	if deps.UserService == nil {
		return nil, errors.NewInternal("invalid user service")
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
	return &loginSvc{
		userSvc:        deps.UserService,
		twoFaSvc:       deps.TwoFAService,
		credentialsSvc: deps.CredentialsService,
		jwtSvc:         deps.JWTService,
		log:            log,
	}, nil
}

func (svc *loginSvc) Login(ctx context.Context, dto *LoginDTO) error {
	// find user
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return errors.WithMessage(ErrPermissionDenied, err.Error())
	}

	// map dto to user
	registeredUser, err := user.MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// if user does not active or not verified return ErrPermissionDenied
	if !registeredUser.IsActive() || !registeredUser.IsVerified {
		return ErrPermissionDenied
	}

	// map from entity to credentials dto
	credentialsDTO := credentials.MapToDTO(registeredUser.Credentials)

	// check password
	if err = svc.credentialsSvc.ValidatePassword(ctx, credentialsDTO, dto.Password); err != nil {
		return err
	}

	return nil
}

func (svc *loginSvc) CheckCode(ctx context.Context, dto *LoginCodeDTO) (*TokenDTO, error) {
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

	// check TwoFA Code
	if err := svc.twoFaSvc.CheckTwoFACode(ctx, dto.Code, *registeredUser.Credentials.SecretOTP); err != nil {
		return nil, err
	}

	// create JWT
	jwtTokenDTO, err := svc.jwtSvc.CreateJWT(ctx, dto.Email)
	if err != nil {
		return nil, errors.WithMessage(ErrUnauthorized, err.Error())
	}
	return &TokenDTO{Token: jwtTokenDTO.Token, ExpireAt: jwtTokenDTO.ExpireAt}, nil
}

// todo: find then delete or diacttivate jwt token
func (svc *loginSvc) Logout(ctx context.Context, email string) error {
	return nil
}
