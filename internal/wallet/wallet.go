package wallet

type Wallet struct {
	WalletName string `bson:"wallet_name"`
	Address    string `bson:"address"`
}
