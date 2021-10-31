package Ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"testing"
)

func TestGenerateEthHDWallet(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections
}

func TestGenerateEthHDWalletAndMakeTransaction(t *testing.T) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	_, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("we have a connection")
	nonce := uint64(0)
	value := big.NewInt(1000000000000000000)
	toAddress := account.Address
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(21000000000)
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	_, err = wallet.SignTx(account, tx, nil)
	if err != nil {
		t.Fatal(err)
	}
}
