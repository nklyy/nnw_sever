package not_working

func init() {
	Coins[LtcType] = NewLTC
	Coins[DogeType] = NewDOGE
	Coins[BtcTestNetType] = NewBtcTestNet
}

type BTCCoins struct {
	*BTC
}

func NewLTC(key *Key) Wallet {
	key.Opt.Params = &LTCParams
	token := NewBTC(key).(*BTC)
	token.Name = "Litecoin"
	token.Symbol = "LTC"

	return &BTCCoins{BTC: token}
}

func NewDOGE(key *Key) Wallet {
	key.Opt.Params = &DOGEParams
	token := NewBTC(key).(*BTC)
	token.Name = "Dogecoin"
	token.Symbol = "DOGE"

	return &BTCCoins{BTC: token}
}

func NewBtcTestNet(key *Key) Wallet {
	key.Opt.Params = &BTCTestnetParams
	token := NewBTC(key).(*BTC)
	token.Name = "BTC Test Net"
	token.Symbol = "BTC"

	return &BTCCoins{BTC: token}
}
