package wallet

import (
	"github.com/google/uuid"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

type Payload struct {
	WalletId string
	Address  string
}

func CreateBTCWallet(backup bool, password, wif, address, mnemonic string) (*Payload, error) {
	//var km *KeyManager
	//var bError error

	//if !backup {
	//	km, bError = NewKeyManager(256, password, "")
	//	if bError != nil {
	//		return nil, bError
	//	}
	//} else {
	//	km, bError = NewKeyManager(256, "", mnemonic)
	//	if bError != nil {
	//		return nil, bError
	//	}
	//}

	//km, err := NewKeyManager(256, "", mnemonic)
	//if err != nil {
	//	return nil, err
	//}

	//masterKey, err := km.GetMasterKey()
	//if err != nil {
	//	return nil, err
	//}

	//passphrase := km.GetPassphrase()
	//if passphrase == "" {
	//	passphrase = "<none>"
	//}
	//
	//key, err := km.GetKey(PurposeBIP44, CoinTypeTestNetBTC, 0, 0, 1)
	//if err != nil {
	//	return nil, err
	//}
	//
	//wif, address, _, _, err := key.Encode(true)
	//if err != nil {
	//	return nil, err
	//}

	walletId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	wallet, err := rpc.CreateWallet(walletId.String())
	if err != nil {
		return nil, err
	}

	err = rpc.EncryptWallet(password, wallet)
	if err != nil {
		return nil, err
	}

	err = rpc.UnLockWallet(password, wallet)
	if err != nil {
		return nil, err
	}

	err = rpc.ImportPrivateKey(wif, walletId.String(), false)
	if err != nil {
		return nil, err
	}

	err = rpc.LockWallet(walletId.String())
	if err != nil {
		return nil, err
	}

	return &Payload{
		WalletId: walletId.String(),
		Address:  address,
	}, nil
}
