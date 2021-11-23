package wallet

type Wallet struct {
	Name       string `bson:"name"`
	WalletName string `bson:"wallet_name"`
	Address    string `bson:"address"`
}

type BTCWallet struct {
	Address    string
	PrivateKey string
}

type ETHWallet struct {
	Address    string
	PrivateKey string
}

var NilWallet *[]*Wallet = nil
