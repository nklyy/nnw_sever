package Bitcoin

import (
	"fmt"
	"testing"
)

func TestTransaction(t *testing.T) {

	//First wallet
	//publicKey1 := "myXDShrmmeMbaTC7boTu1w2PQJJ3DsnTPh"
	privateKey1 := "93NCZ5qRhqSv5MFWPbcUhyHCjxNBoYdRc7oksAwUSq39xWaAH9P"

	//Second wallet
	publicKey2 := "mmNL9hYswxntR5ZJtWafN4PeCJ7cZ4egkj"
	//privateKey2 := "92PezuLo2r1V58mkMEua3uk7QUFU2LD1BHrM4thxom56MRMGUjo"

	rawTx, err := CreateTx(privateKey1,
		publicKey2, 5)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}
