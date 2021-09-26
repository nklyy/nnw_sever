package wallet

// mnemonic language
const (
	English            = "english"
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
)

// zero is default of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	BTC = ZeroQuote + 0
	LTC = ZeroQuote + 2
	ETH = ZeroQuote + 60
	SOL = ZeroQuote + 501

	// btc token
	USDT = BTCToken + 1

	// eth token
	USDC = ETHToken + 2
)

var coinTypes = map[uint32]uint32{
	USDT: BTC,
	USDC: ETH,
}
