package wallet

type Wallet struct {
	Name       string `bson:"name"`
	WalletName string `bson:"wallet_name"`
	Address    string `bson:"address"`
}

var NilWallet *[]*Wallet = nil
