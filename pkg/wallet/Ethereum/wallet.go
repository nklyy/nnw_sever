package Ethereum

import (
	"github.com/ethereum/go-ethereum/crypto"
	"nnw_s/pkg/wallet/Bitcoin/not_working"
)

func init() {
	not_working.Coins[not_working.EthType] = NewETH
}

type ETH struct {
	Name   string
	Symbol string
	Key    *not_working.Key

	// eth token
	contract string
}

func NewETH(key *not_working.Key) not_working.Wallet {
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

func (c *ETH) GetKey() *not_working.Key {
	return c.Key
}

func (c *ETH) GetAddress() (string, error) {
	return crypto.PubkeyToAddress(*c.Key.PublicECDSA).Hex(), nil
}
