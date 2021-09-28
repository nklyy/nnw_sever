package Ethereum

import (
	"github.com/ethereum/go-ethereum/crypto"
	"nnw_s/pkg/wallet"
)

func init() {
	wallet.Coins[wallet.EthType] = NewETH
}

type ETH struct {
	Name   string
	Symbol string
	Key    *wallet.Key

	// eth token
	contract string
}

func NewETH(key *wallet.Key) wallet.Wallet {
	return &ETH{
		Name:   "Ethereum",
		Symbol: "ETH",
		Key:    key,
	}
}

func (c *ETH) GetType() uint32 {
	return c.Key.Opt.CoinType
}

func (c *ETH) GetName() string {
	return c.Name
}

func (c *ETH) GetSymbol() string {
	return c.Symbol
}

func (c *ETH) GetKey() *wallet.Key {
	return c.Key
}

func (c *ETH) GetAddress() (string, error) {
	return crypto.PubkeyToAddress(*c.Key.PublicECDSA).Hex(), nil
}
