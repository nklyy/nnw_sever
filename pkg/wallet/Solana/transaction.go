package Solana

type ISolanaTransaction interface {
	MakeTransfer(fromAccount string, toAccount string, lamports int) error
}

type SolanaTransaction struct {
	SolanaWeb3Client ISolanaWeb3Client
}

func NewSolanaTransaction(solanaWeb3Client ISolanaWeb3Client) ISolanaWallet {
	return &SolanaWallet{
		SolanaWeb3Client: solanaWeb3Client,
	}
}

func (s *SolanaWallet) MakeTransfer(fromAccount string, toAccount string, lamports int) error {
	//secretKey := getSecretKey()
	err := s.SolanaWeb3Client.MakeTransaction(fromAccount, []byte{}, toAccount, lamports)
	return err
}
