package not_working

import (
	"nnw_s/pkg/wallet"
)

func init() {
	wallet.Coins[wallet.BtcType] = NewBTC
}

type BTC struct {
	Name   string
	Symbol string
	Key    *wallet.Key
}

func NewBTC(key *wallet.Key) wallet.Wallet {
	return &BTC{
		Name:   "Bitcoin",
		Symbol: "BTC",
		Key:    key,
	}
}

func (c *BTC) GetType() uint32 {
	return c.Key.Opt.CoinType
}

func (c *BTC) GetName() string {
	return c.Name
}

func (c *BTC) GetSymbol() string {
	return c.Symbol
}

func (c *BTC) GetKey() *wallet.Key {
	return c.Key
}

func (c *BTC) GetAddress() (string, error) {
	return c.Key.AddressBTC()
}
