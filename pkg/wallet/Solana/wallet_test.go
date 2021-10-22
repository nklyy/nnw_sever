package Solana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"net/http"
	"testing"
	"time"
)

type SolanaKeys struct {
	PublicKey string  `json:"public_key"`
	SecretKey []uint8 `json:"secret_key"`
}

func TestGenerateWalletFromMnemonic(t *testing.T) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	resp, err := http.Get("http://localhost:3000/wallet/generate/" + mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var solanaKeys SolanaKeys
	err = json.NewDecoder(resp.Body).Decode(&solanaKeys)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(solanaKeys)
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
