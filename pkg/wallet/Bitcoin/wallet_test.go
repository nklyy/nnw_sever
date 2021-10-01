package Bitcoin

import (
	"fmt"
	"strings"
	"testing"
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
	*/

	km, err := NewKeyManager(128, "", "dwarf unique fork crunch common penalty behind great human gather then usual")
	if err != nil {
		t.Error(err)
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
	fmt.Printf("\t\t\t\t\t\t%s \n\n", "Create Transaction")

	// Transaction
	privWif := "cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD"
	txHash := "9733285947b666571cc826b0243650393b0b9d635a6088736b9b2a7e73a5c862"
	destination := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	amount := int64(40000)
	txFee := int64(10000)
	balance := int64(18300)

	tx, err := CreateTransaction(privWif, txHash, destination, amount, txFee, balance)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(strings.Repeat("-", 106))
	fmt.Printf("%-18s %s\n", "Transaction:", tx)
	// https://live.blockcypher.com/btc-testnet/tx/60f1fb9392121300ecfa1025fe40cb20992ea3fa987ec3c08ef0a6f75f161506/?__cf_chl_jschl_tk__=pmd_d7xmEfC0vuVHctIFgpnwVdUHwDS8Z94xL0I1ZDfiB_U-1633025670-0-gqNtZGzNAmWjcnBszQi9
}
