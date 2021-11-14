package not_working

import (
	"fmt"
	"testing"
)

func TestCrateBtcHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("leader such empower maximum anxiety pilot shadow destroy joke claw correct doctor"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcWallet, _ := master.GetWallet(CoinType(BtcTestNetType), AddressIndex(1))
	btcAddress, _ := btcWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcAddress)

	addressP2WPKH, _ := btcWallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := btcWallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("Bitcoin: ", btcAddress, addressP2WPKH, addressP2WPKHInP2SH)
}

func TestCrateBtcTestHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("tiger rent slam skin fiscal zebra unfold major dune giggle paper axis"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcTestWallet, _ := master.GetWallet(CoinType(BtcTestNetType), AddressIndex(1))
	btcTestAddress, _ := btcTestWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcTestAddress)
}

func TestCrateBtcTestHDWalletAndCreateTransaction(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcTestWallet, _ := master.GetWallet(CoinType(BtcTestNetType), AddressIndex(1))
	btcTestAddress, _ := btcTestWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcTestAddress)
}

func TestGenerateLTCHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ltcWallet, _ := master.GetWallet(CoinType(LtcType))
	ltcAddress, _ := ltcWallet.GetAddress()
	fmt.Println("LTC Address: ", ltcAddress)
}

func TestGenerateDOGEHDWallet(t *testing.T) {
	master, err := NewKey(
		Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	dogeWallet, _ := master.GetWallet(CoinType(DogeType))
	dogeAddress, _ := dogeWallet.GetAddress()
	fmt.Println("DOGE Address: ", dogeAddress)
}
