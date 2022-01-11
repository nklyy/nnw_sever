package helpers

import (
	"math/big"
	"strings"
)

func ConvertHexToDecimal(hex string) *big.Int {
	numberStr := strings.Replace(hex, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	result := new(big.Int)
	result.SetString(numberStr, 16)

	return result
}
