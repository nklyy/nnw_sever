package Ethereum

import (
	"fmt"
	"nnw_s/pkg/wallet"
	"testing"
)

func TestGenerateEthHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ethWallet, _ := master.GetWallet(wallet.CoinType(wallet.EthType))
	ethAddress, _ := ethWallet.GetAddress()
	fmt.Println("Ethereum Address: ", ethAddress)
}

func TestGenerateIOSTHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	IOSTWallet, _ := master.GetWallet(wallet.CoinType(wallet.IOST))
	IOSTAddress, _ := IOSTWallet.GetAddress()
	fmt.Println("IOST Address: ", IOSTAddress)
}

func TestGenerateUSDCHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	USDCWallet, _ := master.GetWallet(wallet.CoinType(wallet.USDC))
	USDCAddress, _ := USDCWallet.GetAddress()
	fmt.Println("USDC Address: ", USDCAddress)
}

func TestGenerateOMGHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	OMGWallet, _ := master.GetWallet(wallet.CoinType(wallet.OMG))
	OMGAddress, _ := OMGWallet.GetAddress()
	fmt.Println("OMG Address: ", OMGAddress)
}

//0x2a5E7ddC6BcC51Ee37FD54C21E5a394DDc48bbf6
