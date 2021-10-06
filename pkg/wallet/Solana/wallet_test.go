package Solana

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"math/big"
	"nnw_s/pkg/wallet"
	"testing"
	"time"
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
