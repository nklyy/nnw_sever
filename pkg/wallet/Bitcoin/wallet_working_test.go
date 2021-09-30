package Bitcoin

import (
	"fmt"
	"strings"
	"testing"
)

func TestWalletAndTransaction(t *testing.T) {
	//fmt.Printf("\n%-34s %-52s %-42s %s\n", "Bitcoin Address", "WIF(Wallet Import Format)", "SegWit(bech32)", "SegWit(nested)")
	//fmt.Println(strings.Repeat("-", 165))
	//
	//wif, address, segwitBech32, segwitNested, err := Generate(true)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Printf("%-34s %s %s %s\n", address, wif, segwitBech32, segwitNested)

	//miy8nyGCWm3jFVD5uJPDPtXw4zDCavdcHS
	//tiger rent slam skin fiscal zebra unfold major dune giggle paper axis
	//indicate drama magic eagle window network jungle stable erode family tuna enact
	//chair column reveal income inside soul blade concert series syrup ivory bulb
	km, err := NewKeyManager(128, "", "tiger rent slam skin fiscal zebra unfold major dune giggle paper axis")
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

	// Transaction
	//cTujmQgVdGYzmZEfhq5gVDpd2EAHF1sZahPmkDnRHmPDEVRYz6eo
	rawTx, err := CreateTx("cReoLGqUKXSuqqJdYeCN7MYDynVQt2zKAZBLG98nhaEuoQiHT4Wt",
		"miimB868qTQ3y8bnwjLUq4Av3e63HZy7nt", 10000)

	if err != nil {
		t.Error(err.Error())
	}
	//
	fmt.Println(strings.Repeat("-", 106))
	fmt.Printf("%-18s %s\n", "Transactio:", rawTx)
	//
	//bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "btc", "test3"}
	//skel, err := bcy.PushTX(rawTx)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Error(err.Error())
	//}
	//fmt.Printf("%+v\n", skel)
	//tx, err := CreateTransaction("cTujmQgVdGYzmZEfhq5gVDpd2EAHF1sZahPmkDnRHmPDEVRYz6eo", "miimB868qTQ3y8bnwjLUq4Av3e63HZy7nt", 1000, "c6950f355835c361dce2e9d6eb511cf56972b67cb34dad5d1fd9f9bc796711a5")
	//fmt.Println("raw signed transaction is: ", tx)
}
