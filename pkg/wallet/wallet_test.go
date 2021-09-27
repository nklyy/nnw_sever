package wallet

import (
	"fmt"
	"testing"
)

func TestCrateBtcHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcWallet, _ := master.GetWallet(CoinType(BTC), AddressIndex(1))
	btcAddress, _ := btcWallet.GetAddress()
	fmt.Println("BTC Address:", btcAddress)

	addressP2WPKH, _ := btcWallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := btcWallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("BTC: ", btcAddress, addressP2WPKH, addressP2WPKHInP2SH)
}

func TestGenerateEthHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ethWallet, _ := master.GetWallet(CoinType(ETH))
	ethAddress, _ := ethWallet.GetAddress()
	fmt.Println("ETH Address: ", ethAddress)
}

func TestGenerateLTCHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ltcWallet, _ := master.GetWallet(CoinType(LTC))
	ltcAddress, _ := ltcWallet.GetAddress()
	fmt.Println("LTC Address: ", ltcAddress)
}

func TestGenerateSOLHDWallet(t *testing.T) {
	//garment inflict make idle duck pepper summer flash target act will access cage charge snow salmon total panic romance foil police hill infant drama
	//Time to hack with only one card: 3830854 years

	//chair column reveal income inside soul blade concert series syrup ivory bulb
	//Time to hack with only one card: 109 seconds
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	solWallet, _ := master.GetWallet(CoinType(SOL))
	solAddress, _ := solWallet.GetAddress()
	fmt.Println("SOL Address: ", solAddress)
}
