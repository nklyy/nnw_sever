package Ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"testing"
)

//func TestGenerateEthHDWallet(t *testing.T) {
//	master, err := not_working.NewKey(
//		not_working.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
//	)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	ethWallet, _ := master.GetWallet(not_working.CoinType(not_working.EthType))
//	ethAddress, _ := ethWallet.GetAddress()
//	fmt.Println("Ethereum Address: ", ethAddress)
//}
//
//func TestGenerateIOSTHDWallet(t *testing.T) {
//	master, err := not_working.NewKey(
//		not_working.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
//	)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	IOSTWallet, _ := master.GetWallet(not_working.CoinType(not_working.IOST))
//	IOSTAddress, _ := IOSTWallet.GetAddress()
//	fmt.Println("IOST Address: ", IOSTAddress)
//}
//
//func TestGenerateUSDCHDWallet(t *testing.T) {
//	master, err := not_working.NewKey(
//		not_working.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
//	)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	USDCWallet, _ := master.GetWallet(not_working.CoinType(not_working.USDC))
//	USDCAddress, _ := USDCWallet.GetAddress()
//	fmt.Println("USDC Address: ", USDCAddress)
//}
//
//func TestGenerateOMGHDWallet(t *testing.T) {
//	master, err := not_working.NewKey(
//		not_working.Mnemonic("chair column reveal income inside soul blade concert series syrup ivory bulb"),
//	)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	OMGWallet, _ := master.GetWallet(not_working.CoinType(not_working.OMG))
//	OMGAddress, _ := OMGWallet.GetAddress()
//	fmt.Println("OMG Address: ", OMGAddress)
//}

//0x2a5E7ddC6BcC51Ee37FD54C21E5a394DDc48bbf6
func TestDecryptKey(t *testing.T) {
	data, err := os.ReadFile("./UTC--2021-12-06T10-54-24.102702324Z--48d029bcebe8c846a9e86b6b715ef9cd526b1641")
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Print(string(data))

	walletData, err := keystore.DecryptKey(data, "asd123")
	if err != nil {
		t.Fatal(err)
	}

	keyBytes := crypto.FromECDSA(walletData.PrivateKey)
	privateKey := hexutil.Encode(keyBytes)
	//fmt.Println(walletData)
	//fmt.Println(privateKey)
	fmt.Printf("%-18s %s\n", "Private key:", privateKey)
}
