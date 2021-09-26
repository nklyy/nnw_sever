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
	ethAddress, _ := ltcWallet.GetAddress()
	fmt.Println("LTC Address: ", ethAddress)
}
