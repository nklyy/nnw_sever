package Ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type IWalletClient interface {
	GetAddressBalance(address string) (string, error)
	CreateWallet(mnemonic string) (*ecdsa.PrivateKey, error)
}

type WalletClient struct {
	client ethclient.Client
}

func NewWalletClient(client ethclient.Client) IWalletClient {
	return &WalletClient{
		client: client,
	}
}

// GetAddressBalance returns the given address balance =P
func (w *WalletClient) GetAddressBalance(address string) (string, error) {
	account := common.HexToAddress(address)
	balance, err := w.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "0", err
	}

	return balance.String(), nil
}

func (w *WalletClient) CreateWallet(mnemonic string) (*ecdsa.PrivateKey, error) {
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
