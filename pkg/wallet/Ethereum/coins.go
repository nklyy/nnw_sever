package Ethereum

import (
	"nnw_s/pkg/wallet/Bitcoin/not_working"
)

func init() {
	not_working.Coins[not_working.USDC] = NewUSDC
	not_working.Coins[not_working.IOST] = NewIOST
	not_working.Coins[not_working.OMG] = NewOMG
}

type ETHCoins struct {
	*ETH
}

func NewIOST(key *not_working.Key) not_working.Wallet {
	token := NewETH(key).(*ETH)
	token.Name = "IOStoken"
	token.Symbol = "IOST"
	token.contract = "0xfa1a856cfa3409cfa145fa4e20eb270df3eb21ab"

	return &ETHCoins{ETH: token}
}

func NewUSDC(key *not_working.Key) not_working.Wallet {
	token := NewETH(key).(*ETH)
	token.Name = "USD Coin"
	token.Symbol = "USDC"
	token.contract = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"

	return &ETHCoins{ETH: token}
}

func NewOMG(key *not_working.Key) not_working.Wallet {
	token := NewETH(key).(*ETH)
	token.Name = "OMG Coin"
	token.Symbol = "OMG"
	token.contract = "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"

	return &ETHCoins{ETH: token}
}
