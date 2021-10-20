package feature

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBuildTransaction(t *testing.T) {
	BuildTransaction()
}

func TestBuildTransactionV2(t *testing.T) {
	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"
	destinationAddress := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	userWalletName := "first"
	userWalletPassword := "password"
	amountToSend := big.NewInt(5000) // amount to send in satoshis (0.01 btc)

	fmt.Printf("%-18s %s\n", "user wallet name:", userWalletName)
	fmt.Printf("%-18s %s\n", "from wallet public address:", fromWalletPublicAddress)
	fmt.Printf("%-18s %s\n", "to wallet public address:", destinationAddress)
	fmt.Printf("%-18s %s\n", "amount:", amountToSend)

	BuildTransactionV2(fromWalletPublicAddress, destinationAddress, userWalletName, userWalletPassword, amountToSend)
}
