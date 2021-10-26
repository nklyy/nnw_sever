package Solana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SolanaAccount struct {
	PublicKey string          `json:"publicKey"`
	SecretKey json.RawMessage `json:"secretKey"`
}

type Transaction struct {
	FromAccount SolanaAccount `json:"fromAccount"`
	ToAccount   SolanaAccount `json:"toAccount"`
	Lamports    int           `json:"lamports"`
}

func CreateWalletFromMnemonic(mnemonic string) (*SolanaAccount, error) {
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

func MakeAirDrop(address string) error {
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

func MakeTransaction(fromAccountPubKey string, fromAccountSecKey []byte, toAccountPubKey string, lamports int) error {
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
