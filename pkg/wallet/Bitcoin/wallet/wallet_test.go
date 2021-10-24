package wallet

import (
	"fmt"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
	"strings"
	"testing"
	"time"
)

func TestWalletAndTransaction(t *testing.T) {
	/*
		fmt.Printf("\n%-34s %-52s %-42s %s\n", "Bitcoin Address", "WIF(Wallet Import Format)", "SegWit(bech32)", "SegWit(nested)")
		fmt.Println(strings.Repeat("-", 165))

		wif, address, segwitBech32, segwitNested, err := Generate(true)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%-34s %s %s %s\n", address, wif, segwitBech32, segwitNested)
	*/

	/*
		***************************************************************************************
		WALLET #1
		BIP39 Mnemonic: 	dwarf unique fork crunch common penalty behind great human gather then usual
		m/44'/1'/0'/0/1:	mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU
		WIF: 				cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD
		***************************************************************************************
		WALLET #2
		BIP39 Mnemonic: 	leader such empower maximum anxiety pilot shadow destroy joke claw correct doctor
		m/44'/1'/0'/0/1:	mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u
		WIF: 				cP4dZeLM1U39DAaui6q4rF2KwMXPJSA67znfJ2Y22CdzbqVMp2mb
		***************************************************************************************
		WALLET #3
		BIP39 Mnemonic: 	gorilla chronic bronze random glass jar deny ten exotic female kind history
		m/44'/1'/0'/0/1:	mvdu6WEXfk75gjcwm8hjSE5kwHMLy9BMfA
		WIF: 				cMg7YGBar4sMMMswRP8EfrdhgHPYaFEmJfPMBJ1jNf5UrMQCv4DH
		***************************************************************************************
		WALLET #3
		BIP39 Mnemonic: 	security cinnamon absent that side muscle pigeon fat habit sadness veteran subject
		m/44'/1'/0'/0/1:	mrvZjXUNupoEpQf6KsgiVgzLz7DUca6Kfv
		WIF: 				cTcEm5gFJ8H9FCxkuQLdBdSaxx7TPLGkZcE2JjWfZRKK1DcwdwMF
		***************************************************************************************
	*/

	// 128: 12 phrases
	// 256: 24 phrases

	var km *KeyManager
	var bError error
	backUp := false

	if backUp {
		km, bError = NewKeyManager(128, "", "")
		if bError != nil {
			t.Error(bError)
		}
	} else {
		km, bError = NewKeyManager(128, "", "dwarf unique fork crunch common penalty behind great human gather then usual")
		if bError != nil {
			t.Error(bError)
		}
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		t.Error(err)
	}
	passphrase := km.GetPassphrase()
	if passphrase == "" {
		passphrase = "<none>"
	}
	fmt.Printf("\n%-18s %s\n", "BIP39 Mnemonic:", km.GetMnemonic())
	fmt.Printf("%-18s %s\n", "BIP39 Passphrase:", passphrase)
	fmt.Printf("%-18s %x\n", "BIP39 Seed:", km.GetSeed())
	fmt.Printf("%-18s %s\n", "BIP32 Public:", masterKey.PublicKey().B58Serialize())
	fmt.Printf("%-18s %s\n", "BIP32 Root Key:", masterKey.B58Serialize())

	fmt.Printf("\n%-18s %-34s %-52s\n", "Path(BIP44)", "Bitcoin Address", "WIF(Wallet Import Format)")
	fmt.Println(strings.Repeat("-", 106))

	key, err := km.GetKey(PurposeBIP44, CoinTypeTestNetBTC, 0, 0, 1)
	if err != nil {
		t.Error(err)
	}
	wif, address, _, _, err := key.Encode(true)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%-18s %-34s %s\n", key.GetPath(), address, wif)
	fmt.Println(strings.Repeat("-", 106))

	createdWalletName, err := rpc.CreateWallet("ninth")
	if err != nil {
		t.Error(err)
	}

	err = rpc.EncryptWallet("password", createdWalletName)
	if err != nil {
		t.Error(err)
	}

	err = rpc.UnLockWallet("password", createdWalletName)
	if err != nil {
		t.Error(err)
	}

	if backUp {
		go func() {
			err := rpc.ImportPrivateKey(wif, createdWalletName, true)
			if err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(2 * time.Second)
	} else {
		err := rpc.ImportPrivateKey(wif, createdWalletName, false)
		if err != nil {
			t.Error(err)
		}
	}

	fmt.Printf("%-18s %s\n", "Your wallet:", address)
	fmt.Printf("%-18s %s\n", "Your mnemonic:", km.GetMnemonic())
	fmt.Println(strings.Repeat("-", 106))
}
