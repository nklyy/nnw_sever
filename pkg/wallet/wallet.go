package wallet

type Wallet struct {
	Name     string `bson:"name"`
	WalletId string `bson:"wallet_id"`
	Address  string `bson:"address"`
}

type BTCWallet struct {
	Address    string
	PrivateKey string
}

var BTCCoinType = uint32(1)

type ETHWallet struct {
	Address    string
	PrivateKey string
}

var ETHCoinType = uint32(60)

var NilWallet *[]*Wallet = nil
