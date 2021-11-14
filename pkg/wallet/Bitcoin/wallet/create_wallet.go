package wallet

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type Payload struct {
	WalletName string
	Address    string
	Mnemonic   string
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

	fmt.Printf("%-18s %-34s %s\n", key.GetPath(), address, wif)
	fmt.Println(strings.Repeat("-", 106))

	walletUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Payload{
		WalletName: walletUuid.String(),
		Address:    address,
		Mnemonic:   km.GetMnemonic(),
	}, nil
}
