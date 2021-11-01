package Ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"testing"
)

func TestGenerateEthHDWalletAndMakeTransaction(t *testing.T) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	privateKeyECDSA, err := createHdWallet(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	gasTip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		t.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		t.Fatal(err)
	}

	nonce, err := client.NonceAt(ctx, crypto.PubkeyToAddress(privateKeyECDSA.PublicKey), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatal(err)
	}

	value := big.NewInt(1000000000000000000)
	toAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	gasLimit := uint64(21000)

	fmt.Println("toAddress:", &toAddress)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: gasTip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      []byte{},
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		t.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
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

	return privateKeyECDSA, nil
}
