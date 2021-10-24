package btc

import (
	"context"
	"github.com/sirupsen/logrus"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/wallet/Bitcoin/wallet"
)

//go:generate mockgen -source=wallet_service.go -destination=mocks/wallet_service_mock.go
type WalletService interface {
	CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*wallet.Payload, error)
}

type walletSvc struct {
	userSvc        user.Service
	twoFaSvc       twofa.Service
	jwtSvc         jwt.Service
	credentialsSvc credentials.Service

	log *logrus.Logger
}

type ServiceDeps struct {
	UserService        user.Service
	TwoFAService       twofa.Service
	JWTService         jwt.Service
	CredentialsService credentials.Service
}

func NewWalletService(log *logrus.Logger, deps *ServiceDeps) (WalletService, error) {
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
	return &walletSvc{
		userSvc:        deps.UserService,
		twoFaSvc:       deps.TwoFAService,
		jwtSvc:         deps.JWTService,
		credentialsSvc: deps.CredentialsService,
		log:            log,
	}, nil
}

func (svc *walletSvc) CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*wallet.Payload, error) {
	userDTO, _ := svc.userSvc.GetUserByEmail(ctx, email)

	decodePass, err := svc.credentialsSvc.DecodePassword(ctx, dto.Password)
	if err != nil {
		return nil, err
	}

	walletPayload, err := wallet.CreateBTCWallet(*dto.Backup, decodePass, "")
	if err != nil {
		return nil, err
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, err
	}

	userEntity.SetBtcWallet(walletPayload.WalletName)

	// map back to dto
	userDTO = user.MapToDTO(userEntity)

	// update user entity in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user's status field: %v", err)
		return nil, err
	}

	return walletPayload, nil
}
