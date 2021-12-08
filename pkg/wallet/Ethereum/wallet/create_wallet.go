package wallet

import (
	"errors"
	"github.com/google/uuid"
	"nnw_s/pkg/wallet/Ethereum/rpc"
)

type Payload struct {
	WalletId string
	Address  string
}

func CreateETHWallet(password, privateKey string) (*Payload, error) {
	walletId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	address, err := rpc.ImportPrivateKey(privateKey, password)
	if err != nil {
		return nil, err
	}

	locked, err := rpc.LockWallet(address)
	if err != nil {
		return nil, err
	}

	if !locked {
		return nil, errors.New("Wallet doesn't lock. ")
	}

	return &Payload{
		WalletId: walletId.String(),
		Address:  address,
	}, nil
}
