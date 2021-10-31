package Solana

import (
	"fmt"
	"github.com/google/uuid"
)

type ISolanaWallet interface {
	CreateWallet(mnemonic string) (*Payload, error)
	GetBalance(pubKey string) string
}

type Payload struct {
	WalletName string
	Address    string
	Mnemonic   string
}

type SolanaWallet struct {
	SolanaWeb3Client ISolanaWeb3Client
}

func NewSolanaWallet(solanaWeb3Client ISolanaWeb3Client) ISolanaWallet {
	return &SolanaWallet{
		SolanaWeb3Client: solanaWeb3Client,
	}
}

func (s *SolanaWallet) CreateWallet(mnemonic string) (*Payload, error) {
	account, err := s.SolanaWeb3Client.CreateWalletFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	fmt.Println(account.PublicKey)
	fmt.Println(account.SecretKey)

	walletUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Payload{
		WalletName: walletUuid.String(),
		Address:    account.PublicKey,
		Mnemonic:   mnemonic,
	}, nil
}

func (s *SolanaWallet) GetBalance(pubKey string) string {
	return ""
}
