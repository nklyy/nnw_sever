package wallet

func init() {
	Coins[SOL] = newSolana
}

type solana struct {
	name   string
	symbol string
	key    *Key

	// eth token
	contract string
}

func newSolana(key *Key) Wallet {
	return &solana{
		name:   "Solana",
		symbol: "SOL",
		key:    key,
	}
}

func (c *solana) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *solana) GetName() string {
	return c.name
}

func (c *solana) GetSymbol() string {
	return c.symbol
}

func (c *solana) GetKey() *Key {
	return c.key
}

func (c *solana) GetAddress() (string, error) {
	//return crypto.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
	return "", nil
}
