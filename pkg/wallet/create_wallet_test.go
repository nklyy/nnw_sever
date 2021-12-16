package wallet

import (
	"fmt"
	"math/big"
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
	fmt.Printf("%-18s %s\n", "ETH private key:", ethWallet.PrivateKey[2:])
	fmt.Println(strings.Repeat("-", 106))
}

func TestHexadecimalToDecimal(t *testing.T) {
	numberStr := strings.Replace("0x768d0e3d72f76117e9e", "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	fmt.Println(numberStr)

	i := new(big.Int)
	i.SetString(numberStr, 16)
	fmt.Println(i)
	fmt.Printf("Output float %v.6\n", float64(i.Int64())/1e-18)

	//num, err := strconv.ParseUint("214cee87a5b0eb10b97ea", 16, 64)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(num)

	//output, err := strconv.ParseInt(numberStr, 16, 64)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Printf("Output %d\n", output)
}
