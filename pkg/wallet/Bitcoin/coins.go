package Bitcoin

import (
	"nnw_s/pkg/wallet"
	"nnw_s/pkg/wallet/Bitcoin/not_working"
)

func init() {
	wallet.Coins[wallet.LtcType] = NewLTC
	wallet.Coins[wallet.DogeType] = NewDOGE
	wallet.Coins[wallet.BtcTestNetType] = NewBtcTestNet
}

type BTCCoins struct {
	*not_working.BTC
}

func NewLTC(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &LTCParams
	token := not_working.NewBTC(key).(*not_working.BTC)
	token.Name = "Litecoin"
	token.Symbol = "LTC"

	return &BTCCoins{BTC: token}
}

func NewDOGE(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &DOGEParams
	token := not_working.NewBTC(key).(*not_working.BTC)
	token.Name = "Dogecoin"
	token.Symbol = "DOGE"

	return &BTCCoins{BTC: token}
}

func NewBtcTestNet(key *wallet.Key) wallet.Wallet {
	key.Opt.Params = &BTCTestnetParams
	token := not_working.NewBTC(key).(*not_working.BTC)
	token.Name = "BTC Test Net"
	token.Symbol = "BTC"

	return &BTCCoins{BTC: token}
}
