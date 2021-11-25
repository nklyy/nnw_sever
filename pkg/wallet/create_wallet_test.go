package wallet

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	btcKey, err := CreateWallet(BTCCoinType, "pepper fitness kangaroo awesome planet cave melt tide vote wing ramp trim connect estate ball add language absorb web cotton choice roast fluid guess")
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

	ethKey, err := CreateWallet(ETHCoinType, "pepper fitness kangaroo awesome planet cave melt tide vote wing ramp trim connect estate ball add language absorb web cotton choice roast fluid guess")
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

func TestHexadecimalToDecimal(t *testing.T) {
	numberStr := strings.Replace("0x72966d", "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	output, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Output %d\n", output)
}
