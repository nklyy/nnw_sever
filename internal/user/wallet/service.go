package wallet

import (
	"context"
	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"
	"nnw_s/pkg/wallet"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
	btc_transaction "nnw_s/pkg/wallet/Bitcoin/transaction"
	btc_wallet "nnw_s/pkg/wallet/Bitcoin/wallet"
)

//go:generate mockgen -source=wallet_service.go -destination=mocks/wallet_service_mock.go
type Service interface {
	CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*string, error)
	GetWallet(ctx context.Context, email string, walletId string) (*wallet.Wallet, error)
	GetBalance(ctx context.Context, dto *GetWalletBalanceDTO) (*BalanceDTO, error)
	GetWalletTx(ctx context.Context, dto *GetWalletTxDTO) ([]*TxsDTO, error)

	CreateTx(ctx context.Context, dto *CreateTxDTO) (string, string, error)
	SendTx(ctx context.Context, dto *SendTxDTO, email string) (string, error)
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

func (svc *walletSvc) CreateWallet(ctx context.Context, dto *CreateWalletDTO, email string, shift int) (*string, error) {
	userDTO, _ := svc.userSvc.GetUserByEmail(ctx, email)

	decodePass, err := svc.credentialsSvc.DecodePassword(ctx, dto.Password)
	if err != nil {
		return nil, err
	}

	walletNameMap := []string{"BTC"}
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

	for _, w := range walletNameMap {
		switch w {
		case "BTC":
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
		}
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

func (svc *walletSvc) GetBalance(ctx context.Context, dto *GetWalletBalanceDTO) (*BalanceDTO, error) {

	var balanceInt *big.Int
	var balanceStr *btcutil.Amount

	switch dto.Name {
	case "BTC":
		warning, err := rpc.LoadWallet(dto.WalletId)
		if err != nil {
			return nil, err
		}

		if warning != "" {
			return nil, errors.NewInternal(warning)
		}

		balanceInt, balanceStr, err = rpc.GetBalance(dto.WalletId)
		if err != nil {
			return nil, err
		}
	}

	return &BalanceDTO{
		Balance:    float64(balanceStr.MulF64(1e-8)),
		BalanceInt: balanceInt,
		BalanceStr: balanceStr.String(),
	}, nil
}

func (svc *walletSvc) GetWalletTx(ctx context.Context, dto *GetWalletTxDTO) ([]*TxsDTO, error) {

	var resultTxs []*TxsDTO

	switch dto.Name {
	case "BTC":
		warning, err := rpc.LoadWallet(dto.WalletId)
		if err != nil {
			return nil, err
		}

		if warning != "" {
			return nil, errors.NewInternal(warning)
		}

		_, txs, err := rpc.TransactionList(dto.WalletId)
		if err != nil {
			return nil, err
		}

		//for _, tx := range txs {
		//	arrTxs = append(arrTxs, &TxsDTO{
		//		Address:  tx.Address,
		//		Category: tx.Category,
		//		Amount:   tx.Amount,
		//		Txid:     tx.Txid,
		//	})
		//}

		var sortedTx []string
		for _, tx := range txs {
			if !helpers.ContainsStr(sortedTx, tx.Txid) {
				sortedTx = append(sortedTx, tx.Txid)
			}
		}

		for _, tx := range sortedTx {
			var inputTx []*InputTxDTO
			var outputTx []*OutTxDTO

			rt, err := rpc.GetRawTransaction(dto.WalletId, dto.Address, tx)
			if err != nil {
				return nil, err
			}

			for _, in := range rt.Vin {
				rt, err := rpc.GetRawTransaction(dto.WalletId, dto.Address, in.Txid)
				if err != nil {
					return nil, err
				}

				for _, out := range rt.Vout {
					if out.N == in.Vout {
						inputTx = append(inputTx, &InputTxDTO{
							Address: out.ScriptPubKey.Addresses,
							Value:   out.Value,
						})
					}
				}
			}

			for _, out := range rt.Vout {
				outputTx = append(outputTx, &OutTxDTO{
					Address: out.ScriptPubKey.Addresses,
					Value:   out.Value,
				})
			}

			resultTxs = append(resultTxs, &TxsDTO{
				Txid:          tx,
				Time:          rt.Time,
				Confirmations: rt.Confirmations,
				Input:         inputTx,
				Output:        outputTx,
			})
		}
	}

	return resultTxs, nil
}

func (svc *walletSvc) CreateTx(ctx context.Context, dto *CreateTxDTO) (string, string, error) {

	var notSignTx string
	var fee string

	switch dto.Name {
	case "BTC":
		amount, err := btcutil.NewAmount(dto.Amount)
		if err != nil {
			return "", "", err
		}

		nstx, f, err := btc_transaction.CreateNotSignTx(dto.FromAddress, dto.ToAddress, dto.WalletId, big.NewInt(int64(amount)))
		if err != nil {
			return "", "", err
		}

		notSignTx = nstx
		feeAmount, _ := btcutil.NewAmount(float64(f.Int64()) / 1e8)
		fee = feeAmount.String()
	}

	return notSignTx, fee, nil
}

func (svc *walletSvc) SendTx(ctx context.Context, dto *SendTxDTO, email string) (string, error) {
	userDTO, err := svc.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = svc.twoFaSvc.CheckTwoFACode(ctx, dto.TwoFaCode, userDTO.SecretOTP)
	if err != nil {
		return "", err
	}

	decodePass, err := svc.credentialsSvc.DecodePassword(ctx, dto.Password)
	if err != nil {
		return "", err
	}

	var txHash string

	switch dto.Name {
	case "BTC":
		amount, err := btcutil.NewAmount(dto.Amount)
		if err != nil {
			return "", err
		}

		h, err := btc_transaction.SignAndSendTx(decodePass, dto.WalletId, dto.FromAddress, dto.NotSignTx, big.NewInt(int64(amount)))
		if err != nil {
			return "", err
		}

		txHash = h
	}

	return txHash, nil
}
