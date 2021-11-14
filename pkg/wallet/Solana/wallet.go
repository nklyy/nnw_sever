package Solana

import (
	"fmt"
	"nnw_s/pkg/wallet/Bitcoin/not_working"
)

func init() {
	not_working.Coins[not_working.SolType] = NewSolana
}

type Solana struct {
	Name   string
	Symbol string
	Key    *not_working.Key

	// eth token
	contract string
}

func NewSolana(key *not_working.Key) not_working.Wallet {
	return &Solana{
		Name:   "Solana",
		Symbol: "Solana",
		Key:    key,
	}
}

func (c *Solana) GetType() uint32 {
	return c.Key.Opt.CoinType
}

func (c *Solana) GetName() string {
	return c.Name
}

func (c *Solana) GetSymbol() string {
	return c.Symbol
}

func (c *Solana) GetKey() *not_working.Key {
	return c.Key
}

func (c *Solana) GetAddress() (string, error) {
	fmt.Println(*c.Key.Public.ToECDSA())
	//return crypto.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
	return "", nil
}
