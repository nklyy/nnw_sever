package Ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip39"

	"testing"
)

func TestGenerateEthHDWalletAndMakeTransaction(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		t.Fatal(err)
	}

	walletClient := NewWalletClient(*client)
	transactionClient := NewTransactionClient(*client)

	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("\n" + mnemonic)

	privateKeyECDSA, err := walletClient.CreateWallet(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	fromPrivKey := "c30867c0109fe66e1eb6500e40c97a47aa2e256186049330dfb25698dd353926"
	toAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	fmt.Println("\n toAddress:", &toAddress)

	tx, err := transactionClient.TransferEth(fromPrivKey, toAddress.String(), 100000)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n transaction:", tx)

	balance, err := walletClient.GetAddressBalance(toAddress.String())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n balance:", balance)

	log, err := transactionClient.GetLogTransactions([]common.Address{toAddress})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n log:", log)
}
