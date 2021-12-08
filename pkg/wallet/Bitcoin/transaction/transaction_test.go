package transaction

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestBuildTransaction(t *testing.T) {
	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"
	destinationAddress := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	userWalletId := "first"
	userWalletPassword := "password"
	amountToSend := big.NewInt(5555) // amount to send in satoshis (0.01 btc)

	BuildTransaction(fromWalletPublicAddress, destinationAddress, userWalletId, userWalletPassword, amountToSend)
}

func TestBuildTransactionV2(t *testing.T) {
	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"
	destinationAddress := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	userWalletId := "first"
	userWalletPassword := "password"
	amountToSend := big.NewInt(5555) // amount to send in satoshis (0.01 btc)

	fmt.Printf("%-18s %s\n", "user wallet name:", userWalletId)
	fmt.Printf("%-18s %s\n", "user wallet password:", userWalletPassword)
	fmt.Printf("%-18s %s\n", "from wallet public address:", fromWalletPublicAddress)
	fmt.Printf("%-18s %s\n", "to wallet public address:", destinationAddress)
	fmt.Printf("%-18s %s\n", "amount:", amountToSend)
	fmt.Println(strings.Repeat("-", 106))

	BuildTransactionV2(fromWalletPublicAddress, destinationAddress, userWalletId, userWalletPassword, amountToSend)
}
