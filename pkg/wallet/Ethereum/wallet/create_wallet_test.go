package wallet

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"nnw_s/pkg/wallet/Ethereum/rpc"
	"strconv"
	"testing"
)

func TestCreateETHWallet(t *testing.T) {
	address, err := rpc.ImportPrivateKey("67d0fc18baac0fa03451ccd108c451119929ab6e7467665965b7117fa127896c", "asd")
	if err != nil {
		t.Fatal(err)
	}

	locked, err := rpc.LockWallet(address)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("LOCKED:", locked)

	if !locked {
		t.Fatal(errors.New("Wallet doesn't lock. "))
	}

	fmt.Println("ADDRESS:", address)
}

func TestGetBalance(t *testing.T) {
	balance, err := rpc.GetBalance("0x090b60A52B3789924AfB9ABFBd88f31526DE868D")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(balance)
	fmt.Println(new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(params.Ether)))
}

func TestGetTransactionOfUserWallet(t *testing.T) {
	blockNumberInt, err := rpc.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	startBlock := blockNumberInt.Int64() - 100
	endBlock := blockNumberInt.Int64()
	fmt.Println(startBlock, endBlock)
	var resTx []rpc.Transaction

	for i := startBlock; i < endBlock; i++ {
		hexNum := strconv.FormatInt(i, 16)

		block, err := rpc.GetBlock("0x"+hexNum, true)
		if err != nil {
			t.Fatal(err)
		}

		for _, tx := range block.Transactions {
			if tx.To == "0x090b60a52b3789924afb9abfbd88f31526de868d" || tx.From == "0x090b60a52b3789924afb9abfbd88f31526de868d" {
				resTx = append(resTx, tx)
			}
		}
	}

	fmt.Println(resTx)
}

func TestSendTransaction(t *testing.T) {
	unlocked, err := rpc.UnlockWallet("0x090b60A52B3789924AfB9ABFBd88f31526DE868D", "asd123")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UNLOCKED:", unlocked)

	if !unlocked {
		t.Fatal(errors.New("Wallet doesn't unlock "))
	}

	//0x214e8348c4f0000 = 0.15 eth
	txHash, err := rpc.SendTransaction("0x090b60a52b3789924afb9abfbd88f31526de868d", "0x0c8bcef257c70727cd5016318b38ce2c63669437", "0x214e8348c4f0000")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Tx HASH: ", txHash)

	locked, err := rpc.LockWallet("0x090b60a52b3789924afb9abfbd88f31526de868d")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("LOCKED:", locked)

	if !locked {
		t.Fatal(errors.New("Wallet doesn't lock. "))
	}
}
