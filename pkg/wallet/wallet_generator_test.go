package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"math/big"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestGenerateWallet(t *testing.T) {
	secretPassphrase := "secret"
	generated := Generate(secretPassphrase)

	// Display wallet and keys
	fmt.Println("Mnemonic: ", generated.Mnemonic)
	fmt.Println("Master private key: ", generated.MasterKey)
	fmt.Println("Master public key: ", generated.PublicKey)

	mnemonicEncrypted := Encrypt([]byte(generated.Mnemonic), secretPassphrase)
	mnemonicDecrypted := Decrypt(mnemonicEncrypted, secretPassphrase)

	fmt.Println("Decrypted Mnemonic:", string(mnemonicDecrypted))

	if string(mnemonicDecrypted) != generated.Mnemonic {
		t.Error("mnemonic does not valid!")
	}

}

func TestGenerateSolWalletMakeTransaction(t *testing.T) {
	// Create a new account:
	account := solana.NewWallet()
	fmt.Println("account private key:", account.PrivateKey)
	fmt.Println("account public key:", account.PublicKey())

	// Create a new RPC client:
	client := rpc.New(rpc.TestNet_RPC)
	{
		// Airdrop 100 SOL to the new account:
		out, err := client.RequestAirdrop(
			context.TODO(),
			account.PublicKey(),
			solana.LAMPORTS_PER_SOL*100,
			rpc.CommitmentFinalized,
		)

		if err != nil {
			t.Error(err.Error())
		}
		fmt.Println(out.String())
	}

	{
		out, err := client.GetBalance(
			context.TODO(),
			account.PublicKey(),
			rpc.CommitmentFinalized,
		)
		if err != nil {
			panic(err)
		}
		spew.Dump(out)
		spew.Dump(out.Value) // total lamports on the account; 1 sol = 1000000000 lamports

		var lamportsOnAccount = new(big.Float).SetUint64(uint64(out.Value))
		// Convert lamports to sol:
		var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))

		// WARNING: this is not a precise conversion.
		fmt.Println("◎", solBalance.Text('f', 10))
	}

}

// 1) install npm install -g ganache-cli
// 2) run in console -> ganache-cli
func TestETHGetBalanceByAccount(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Error(err.Error())
	}
	account := common.HexToAddress("0xA0415a9C644BF7E2BDB181adFCd3ba6144331507")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(balance)
}

//deterministic wallet
func TestETHCreateWallet(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 9a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
