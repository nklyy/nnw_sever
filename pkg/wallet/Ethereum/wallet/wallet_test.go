package wallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"testing"
)

func TestDecryptKey(t *testing.T) {
	data, err := os.ReadFile("./UTC--2021-12-06T10-54-24.102702324Z--48d029bcebe8c846a9e86b6b715ef9cd526b1641")
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Print(string(data))

	walletData, err := keystore.DecryptKey(data, "asd123")
	if err != nil {
		t.Fatal(err)
	}

	keyBytes := crypto.FromECDSA(walletData.PrivateKey)
	privateKey := hexutil.Encode(keyBytes)
	//fmt.Println(walletData)
	//fmt.Println(privateKey)
	fmt.Printf("%-18s %s\n", "Private key:", privateKey)
}
