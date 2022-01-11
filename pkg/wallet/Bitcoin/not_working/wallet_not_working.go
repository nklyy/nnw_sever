package not_working

func init() {
	Coins[BtcType] = NewBTC
}

type BTC struct {
	Name   string
	Symbol string
	Key    *Key
}

func NewBTC(key *Key) Wallet {
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

func (c *BTC) GetKey() *Key {
	return c.Key
}

func (c *BTC) GetAddress() (string, error) {
	return c.Key.AddressBTC()
}
