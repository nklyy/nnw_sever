package wallet

import (
	"github.com/google/uuid"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

type Payload struct {
	WalletName string
	Address    string
}

func CreateBTCWallet(backup bool, password, mnemonic string) (*Payload, error) {
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

	km, err := NewKeyManager(256, "", mnemonic)
	if err != nil {
		return nil, err
	}

	//masterKey, err := km.GetMasterKey()
	//if err != nil {
	//	return nil, err
	//}

	passphrase := km.GetPassphrase()
	if passphrase == "" {
		passphrase = "<none>"
	}

	key, err := km.GetKey(PurposeBIP44, CoinTypeTestNetBTC, 0, 0, 1)
	if err != nil {
		return nil, err
	}

	wif, address, _, _, err := key.Encode(true)
	if err != nil {
		return nil, err
	}

	walletUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	createdWalletName, err := rpc.CreateWallet(walletUuid.String())
	if err != nil {
		return nil, err
	}

	err = rpc.EncryptWallet(password, createdWalletName)
	if err != nil {
		return nil, err
	}

	err = rpc.UnLockWallet(password, createdWalletName)
	if err != nil {
		return nil, err
	}

	err = rpc.ImportPrivateKey(wif, walletUuid.String(), false)
	if err != nil {
		return nil, err
	}

	err = rpc.LockWallet(walletUuid.String())
	if err != nil {
		return nil, err
	}

	return &Payload{
		WalletName: walletUuid.String(),
		Address:    address,
	}, nil
}
