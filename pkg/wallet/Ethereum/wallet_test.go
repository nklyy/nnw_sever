package Ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"nnw_s/pkg/wallet/Ethereum/modules"
	"testing"
)

func TestGenerateEthHDWalletAndMakeTransaction(t *testing.T) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("\n" + mnemonic)

	privateKeyECDSA, err := createHdWallet(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		t.Fatal(err)
	}

	fromPrivKey := "2f5199067115f233332df47de048fe75e0df518b34ee0dafd4cfb5bc44cd212f"
	toAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	fmt.Println("\n toAddress:", &toAddress)

	tx, err := modules.TransferEth(*client, fromPrivKey, toAddress.String(), 100000)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n transaction:", tx)

	balance, err := modules.GetAddressBalance(*client, toAddress.String())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n balance:", balance)

	log, err := modules.GetLogTransactions(*client, []common.Address{toAddress})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n log:", log)
}

func createHdWallet(mnemonic string) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Account address: %s\n", account.Address.Hex())
	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return nil, err
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	fmt.Print("Private key:", privateKey)

	return privateKeyECDSA, nil
}
