package Solana

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"testing"
	"time"
)

func TestCreateWalletsAndMakeTransaction(t *testing.T) {
	entropy1, _ := bip39.NewEntropy(256)
	mnemonic1, _ := bip39.NewMnemonic(entropy1)
	entropy2, _ := bip39.NewEntropy(256)
	mnemonic2, _ := bip39.NewMnemonic(entropy2)

	firstWallet, err := CreateWalletFromMnemonic(mnemonic1)
	if err != nil {
		t.Error(err)
	}

	secondWallet, err := CreateWalletFromMnemonic(mnemonic2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("first wallet: ", firstWallet)
	fmt.Println("second wallet: ", secondWallet)
	fmt.Println(string(firstWallet.SecretKey))
	err = MakeAirDrop(firstWallet.PublicKey)
	if err != nil {
		t.Error(err)
	}

	err = MakeTransaction(firstWallet.PublicKey, firstWallet.SecretKey, secondWallet.PublicKey, 100)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateDeterministicSolanaWalletAndMakeAirDropTransaction(t *testing.T) {
	// Create a new account:
	account := solana.NewWallet()
	fmt.Println("account private key:", account.PrivateKey)
	fmt.Println("account public key:", account.PublicKey())

	// Create a new RPC client:
	client := rpc.New(rpc.TestNet_RPC)

	// Airdrop 1 SOL to the new account:
	out, err := client.RequestAirdrop(
		context.TODO(),
		account.PublicKey(),
		solana.LAMPORTS_PER_SOL*1,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("airdrop transaction signature:", out)

	fmt.Println("waiting for transaction executed!")
	time.Sleep(time.Second * 30)

	//get balance
	GetBalance(account.PublicKey())
}

func GetBalance(pubKey solana.PublicKey) {
	endpoint := rpc.TestNet_RPC
	client := rpc.New(endpoint)

	out, err := client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		fmt.Println(err)
	}

	var lamportsOnAccount = new(big.Float).SetUint64(uint64(out.Value))
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))

	// WARNING: this is not a precise conversion.
	fmt.Println("◎", solBalance.Text('f', 10))
}
