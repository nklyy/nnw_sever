package Bitcoin

import (
	"fmt"
	"strings"
	"testing"
)

func TestWallet_V2(t *testing.T) {
	//fmt.Printf("\n%-34s %-52s %-42s %s\n", "Bitcoin Address", "WIF(Wallet Import Format)", "SegWit(bech32)", "SegWit(nested)")
	//fmt.Println(strings.Repeat("-", 165))
	//
	//wif, address, segwitBech32, segwitNested, err := Generate(true)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Printf("%-34s %s %s %s\n", address, wif, segwitBech32, segwitNested)

	//miy8nyGCWm3jFVD5uJPDPtXw4zDCavdcHS
	km, err := NewKeyManager(128, "", "chair column reveal income inside soul blade concert series syrup ivory bulb")
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

	//mrR75k6L2vHafqamjhj7535XTFGwwsWwdK cPe1P7WzHyZpQURUGP2QMrpdyvCCWckQX1auhQpVeZbLtSDee8pL
	fmt.Printf("%-18s %-34s %s\n", key.GetPath(), address, wif)
}
