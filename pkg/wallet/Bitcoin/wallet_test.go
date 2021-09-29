package Bitcoin

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"nnw_s/pkg/wallet"
	"testing"
)

func TestCrateBtcHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcWallet, _ := master.GetWallet(wallet.CoinType(wallet.BtcType), wallet.AddressIndex(1))
	btcAddress, _ := btcWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcAddress)

	addressP2WPKH, _ := btcWallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := btcWallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("Bitcoin: ", btcAddress, addressP2WPKH, addressP2WPKHInP2SH)
}

func TestCrateBtcTestHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("tiger rent slam skin fiscal zebra unfold major dune giggle paper axis"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcTestWallet, _ := master.GetWallet(wallet.CoinType(wallet.BtcTestNetType), wallet.AddressIndex(1))
	btcTestAddress, _ := btcTestWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcTestAddress)
}

func TestCrateBtcTestHDWalletAndCreateTransaction(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	btcTestWallet, _ := master.GetWallet(wallet.CoinType(wallet.BtcTestNetType), wallet.AddressIndex(1))
	btcTestAddress, _ := btcTestWallet.GetAddress()
	fmt.Println("Bitcoin Address:", btcTestAddress)

	wif, err := btcutil.DecodeWIF("cTujmQgVdGYzmZEfhq5gVDpd2EAHF1sZahPmkDnRHmPDEVRYz6eo")
	if err != nil {
		t.Error(err.Error())
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(addrPubKey.EncodeAddress())

	//addresspubkey, _ := btcutil.NewAddressPubKey(master.Public.SerializeUncompressed(), &chaincfg.TestNet3Params)
	//fmt.Println("ASDASDADASDA", addresspubkey.EncodeAddress())

	//wif, err := master.PrivateWIF(false)
	//if err != nil {
	//	t.Error(err.Error())
	//}
	//fmt.Println(wif)

	//93U3gZk1gDbJvpWuLji9Ce4XfbLTyQHMqRyZ7Pp2HG3sSVC1E4M
	//mx3xkQ4EHWYmAcg5N26J4cCsHHg8rDaqod
	//rawTx, err := CreateTx("93U3gZk1gDbJvpWuLji9Ce4XfbLTyQHMqRyZ7Pp2HG3sSVC1E4M",
	//	"mx3xkQ4EHWYmAcg5N26J4cCsHHg8rDaqod", 10000)
	//
	//if err != nil {
	//	t.Error(err.Error())
	//}
	//
	//fmt.Println("raw signed transaction is: ", rawTx)

	//tx, err := CreateTransaction("5JVeg3qeHcqhhHLApgt6RLXDum7nejkrLV8DDVAfXni783PcqYp", "miimB868qTQ3y8bnwjLUq4Av3e63HZy7nt", 10000, "c6950f355835c361dce2e9d6eb511cf56972b67cb34dad5d1fd9f9bc796711a5")
	//fmt.Println("raw signed transaction is: ", tx)
	//
	//addresspubkey, _ := btcutil.NewAddressPubKey(master.Private.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	//fmt.Println("ASDASDADASDA", addresspubkey.EncodeAddress())
	////tiger rent slam skin fiscal zebra unfold major dune giggle paper axis
}

func TestGenerateLTCHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ltcWallet, _ := master.GetWallet(wallet.CoinType(wallet.LtcType))
	ltcAddress, _ := ltcWallet.GetAddress()
	fmt.Println("LTC Address: ", ltcAddress)
}

func TestGenerateDOGEHDWallet(t *testing.T) {
	master, err := wallet.NewKey(
		wallet.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	dogeWallet, _ := master.GetWallet(wallet.CoinType(wallet.DogeType))
	dogeAddress, _ := dogeWallet.GetAddress()
	fmt.Println("DOGE Address: ", dogeAddress)
}
