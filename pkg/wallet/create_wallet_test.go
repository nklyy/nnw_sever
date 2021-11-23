package wallet

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	btcCoinType := uint32(2)
	btcKey, err := CreateWallet(btcCoinType, "pepper fitness kangaroo awesome planet cave melt tide vote wing ramp trim connect estate ball add language absorb web cotton choice roast fluid guess")
	if err != nil {
		t.Fatal(err)
	}

	btcWallet, err := ToBTCWallet(btcKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "BTC Address:", btcWallet.Address)
	fmt.Printf("%-18s %s\n", "BTC private key:", btcWallet.PrivateKey)
	fmt.Println(strings.Repeat("-", 106))

	ethCoinType := uint32(60)
	ethKey, err := CreateWallet(ethCoinType, "pepper fitness kangaroo awesome planet cave melt tide vote wing ramp trim connect estate ball add language absorb web cotton choice roast fluid guess")
	if err != nil {
		t.Fatal(err)
	}
	ethWallet, err := ToETHWallet(ethKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "ETH Address:", ethWallet.Address)
	fmt.Printf("%-18s %s\n", "ETH private key:", ethWallet.PrivateKey)
	fmt.Println(strings.Repeat("-", 106))
}
