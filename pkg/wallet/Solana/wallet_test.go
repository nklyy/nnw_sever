package Solana

import (
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestCreateWalletsAndMakeTransaction(t *testing.T) {
	entropy1, _ := bip39.NewEntropy(256)
	mnemonic1, _ := bip39.NewMnemonic(entropy1)
	entropy2, _ := bip39.NewEntropy(256)
	mnemonic2, _ := bip39.NewMnemonic(entropy2)

	web3Client := NewSolanaWeb3Client()

	firstWallet, err := web3Client.CreateWalletFromMnemonic(mnemonic1)
	if err != nil {
		t.Error(err)
	}

	secondWallet, err := web3Client.CreateWalletFromMnemonic(mnemonic2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("first wallet: ", firstWallet)
	fmt.Println("second wallet: ", secondWallet)
	err = web3Client.MakeAirDrop(firstWallet.PublicKey)
	if err != nil {
		t.Error(err)
	}

	err = web3Client.MakeTransaction(firstWallet.PublicKey, firstWallet.SecretKey, secondWallet.PublicKey, 100)
	if err != nil {
		t.Error(err)
	}
}
