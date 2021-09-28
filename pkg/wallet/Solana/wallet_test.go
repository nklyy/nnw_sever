package Solana

import (
	"fmt"
	"nnw_s/pkg/wallet"
	"testing"
)

func TestGenerateSOLHDWallet(t *testing.T) {
	//garment inflict make idle duck pepper summer flash target act will access cage charge snow salmon total panic romance foil police hill infant drama
	//Time to hack with only one card: 3830854 years

	//chair column reveal income inside soul blade concert series syrup ivory bulb
	//Time to hack with only one card: 109 seconds
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	solWallet, _ := master.GetWallet(wallet.CoinType(wallet.SolType))
	solAddress, _ := solWallet.GetAddress()
	fmt.Println("Solana Address: ", solAddress)
}
