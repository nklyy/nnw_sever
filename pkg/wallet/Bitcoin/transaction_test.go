package Bitcoin

import (
	"fmt"
	"nnw_s/pkg/wallet"
	"testing"
)

func createTestWalletByMnemonic(mnemonic string) (wallet.Wallet, string, string, string) {
	master, err := wallet.NewKey(
		wallet.Mnemonic(mnemonic),
	)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	btcWallet, _ := master.GetWallet(wallet.CoinType(wallet.BtcType), wallet.AddressIndex(1))
	btcAddress, _ := btcWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcAddress)

	addressP2WPKH, _ := btcWallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := btcWallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("Bitcoin: ", btcAddress, addressP2WPKH, addressP2WPKHInP2SH)

	return btcWallet, btcAddress, addressP2WPKH, addressP2WPKHInP2SH
}

func TestTransaction(t *testing.T) {

	//Create first wallet
	mnemonic1 := "birth blood link boss join action rib gold night disagree pear gate spoon kit coral approve toe guitar dove fault season arrange script convince"
	_, btcAddress1, addressP2WPKH1, addressP2WPKHInP2SH1 := createTestWalletByMnemonic(mnemonic1)
	fmt.Println(btcAddress1, addressP2WPKH1, addressP2WPKHInP2SH1)
	//Create second wallet
	mnemonic2 := "program harsh crime spot squeeze country cry dizzy bread later inform such success stone misery attract wonder choose stool consider elder uphold oak junior"
	_, btcAddress2, addressP2WPKH2, addressP2WPKHInP2SH2 := createTestWalletByMnemonic(mnemonic2)
	fmt.Println(btcAddress2, addressP2WPKH2, addressP2WPKHInP2SH2)

	rawTx, err := CreateTx(addressP2WPKH1,
		btcAddress2, 5)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}
