package wallet

import (
	"fmt"
	"nnw_s/pkg/wallet/Ethereum/rpc"
	"testing"
)

func TestCreateETHWallet(t *testing.T) {
	//address, err := rpc.ImportPrivateKey("67d0fc18baac0fa03451ccd108c451119929ab6e7467665965b7117fa127896c", "asd")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//locked, err := rpc.LockWallet(address)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Println("LOCKED:", locked)
	//
	//if !locked {
	//	t.Fatal(errors.New("Wallet doesn't lock. "))
	//}
	//
	//fmt.Println("ADDRESS:", address)

	balance, err := rpc.GetBalance("0x407d73d8a49eeb85d32cf465507dd71d507100c1")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Output %d\n", *balance)
	fmt.Printf("Output float %v.6\n", float64(*balance)/1e-18)

	fmt.Println(balance)

	err = rpc.GetTransactionCount("0x407d73d8a49eeb85d32cf465507dd71d507100c1")
	if err != nil {
		t.Fatal(err)
	}
}
