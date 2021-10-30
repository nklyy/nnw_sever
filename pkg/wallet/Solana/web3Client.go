package Solana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ISolanaWeb3Client interface {
	CreateWalletFromMnemonic(mnemonic string) (*SolanaAccount, error)
	MakeAirDrop(address string) error
	MakeTransaction(fromAccountPubKey string, fromAccountSecKey []byte, toAccountPubKey string, lamports int) error
}

type SolanaWeb3Client struct{}

type SolanaAccount struct {
	PublicKey string          `json:"publicKey"`
	SecretKey json.RawMessage `json:"secretKey"`
}

type Transaction struct {
	FromAccount SolanaAccount `json:"fromAccount"`
	ToAccount   SolanaAccount `json:"toAccount"`
	Lamports    int           `json:"lamports"`
}

func NewSolanaWeb3Client() ISolanaWeb3Client {
	return &SolanaWeb3Client{}
}

func (s *SolanaWeb3Client) CreateWalletFromMnemonic(mnemonic string) (*SolanaAccount, error) {
	resp, err := http.Get("http://localhost:3000/api/v1/wallet/generate/" + mnemonic)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var solanaKeys SolanaAccount
	err = json.NewDecoder(resp.Body).Decode(&solanaKeys)
	if err != nil {
		return nil, err
	}

	return &solanaKeys, nil
}

func (s *SolanaWeb3Client) MakeAirDrop(address string) error {
	json_data, err := json.Marshal(map[string]string{"address": address})
	resp, err := http.Post("http://localhost:3000/test/v1/make/air/drop", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("air drop does not created!")
	}

	return nil
}

func (s *SolanaWeb3Client) MakeTransaction(fromAccountPubKey string, fromAccountSecKey []byte, toAccountPubKey string, lamports int) error {
	body := Transaction{FromAccount: SolanaAccount{
		PublicKey: fromAccountPubKey,
		SecretKey: fromAccountSecKey,
	}, ToAccount: SolanaAccount{
		PublicKey: toAccountPubKey,
	}, Lamports: lamports}

	json_data, err := json.Marshal(body)
	resp, err := http.Post("http://localhost:3000/api/v1/make/transaction", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("air drop does not created!")
	}

	return nil
}
