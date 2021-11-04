package wallet

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/wallet"
	btc_wallet "nnw_s/pkg/wallet/Bitcoin/wallet"
)

//go:generate mockgen -source=wallet_service.go -destination=mocks/wallet_service_mock.go
type Service interface {
	CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*string, error)
	GetWallet(ctx context.Context, email string, walletId string) (*wallet.Wallet, error)
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

func NewWalletService(log *logrus.Logger, deps *ServiceDeps) (Service, error) {
	if deps == nil {
		return nil, errors.NewInternal("invalid service dependencies")
	}
	//if deps.UserService == nil {
	//	return nil, errors.NewInternal("invalid user service")
	//}
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

func (svc *walletSvc) CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*string, error) {
	userDTO, _ := svc.userSvc.GetUserByEmail(ctx, email)

	decodePass, err := svc.credentialsSvc.DecodePassword(ctx, dto.Password)
	if err != nil {
		return nil, err
	}

	// Create BTC wallets
	//var walletPayload *btc_wallet.Payload
	//if *dto.Backup {
	//	// need to put user mnemonic
	//	walletPayload, err = btc_wallet.CreateBTCWallet(*dto.Backup, decodePass, "")
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	walletPayload, err = btc_wallet.CreateBTCWallet(*dto.Backup, decodePass, "")
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	walletPayload, err := btc_wallet.CreateBTCWallet(*dto.Backup, decodePass, mnemonic)

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, err
	}

	var wallets []*wallet.Wallet

	wallets = append(wallets, &wallet.Wallet{
		Name:       "BTC",
		WalletName: walletPayload.WalletName,
		Address:    walletPayload.Address,
	})

	userEntity.SetWallet(&wallets)

	// map back to dto
	userDTO = user.MapToDTO(userEntity)

	// update user entity in storage
	if err = svc.userSvc.UpdateUser(ctx, userDTO); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to update user's status field: %v", err)
		return nil, err
	}

	return &mnemonic, nil
}

func (svc *walletSvc) GetWallet(ctx context.Context, email string, walletId string) (*wallet.Wallet, error) {
	userDTO, err := svc.userSvc.GetUserByWalletID(ctx, email, walletId)
	if err != nil {
		return nil, err
	}

	// map dto to user
	userEntity, err := user.MapToEntity(userDTO)
	if err != nil {
		return nil, err
	}

	var walletPayload wallet.Wallet
	for _, w := range *userEntity.Wallet {
		if w.WalletName == walletId {
			walletPayload.Name = w.Name
			walletPayload.WalletName = w.WalletName
			walletPayload.Address = w.Address
			break
		}
	}

	return &walletPayload, nil
}
