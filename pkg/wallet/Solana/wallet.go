package Solana

import (
	"fmt"
)

type ISolanaWallet interface {
	GetAddress() (string, error)
	CreateWallet(mnemonic string) error
	MakeTransfer(fromAccount string, toAccount string, lamports int) error
	GetBalance(pubKey string) string
}

type SolanaWallet struct {
	SolanaWeb3Client ISolanaWeb3Client
}

func NewSolana(solanaWeb3Client ISolanaWeb3Client) ISolanaWallet {
	return &SolanaWallet{
		SolanaWeb3Client: solanaWeb3Client,
	}
}

func (s *SolanaWallet) CreateWallet(mnemonic string) error {
	account, err := s.SolanaWeb3Client.CreateWalletFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	fmt.Println(account.PublicKey)
	fmt.Println(account.SecretKey)

	return nil
}

func (s *SolanaWallet) MakeTransfer(fromAccount string, toAccount string, lamports int) error {
	//secretKey := getSecretKey()
	err := s.SolanaWeb3Client.MakeTransaction(fromAccount, []byte{}, toAccount, lamports)
	return err
}

func (s *SolanaWallet) GetAddress() (string, error) {
	return "", nil
}

func (s *SolanaWallet) GetBalance(pubKey string) string {
	return ""
}
