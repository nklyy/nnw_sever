package Bitcoin

import (
	"nnw_s/pkg/wallet"
)

func init() {
	wallet.Coins[wallet.LtcType] = NewLTC
	wallet.Coins[wallet.DogeType] = NewDOGE
	wallet.Coins[wallet.BtcTestNetType] = NewBtcTestNet
}

type BTCCoins struct {
	*BTC
}

func NewLTC(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &LTCParams
	token := NewBTC(key).(*BTC)
	token.Name = "Litecoin"
	token.Symbol = "LTC"

	return &BTCCoins{BTC: token}
}

func NewDOGE(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &DOGEParams
	token := NewBTC(key).(*BTC)
	token.Name = "Dogecoin"
	token.Symbol = "DOGE"

	return &BTCCoins{BTC: token}
}

func NewBtcTestNet(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &BTCTestnetParams
	token := NewBTC(key).(*BTC)
	token.Name = "BTC Test Net"
	token.Symbol = "BTC"

	return &BTCCoins{BTC: token}
}
