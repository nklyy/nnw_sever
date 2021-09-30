package Bitcoin

import (
	"fmt"
	"github.com/blockcypher/gobcy"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	privWif := "cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD"
	txHash := "d30777a4c097757d640bac82ef2eb9ff4e089e65e3080e63a873ce49eff53c28"
	destination := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	amount := int64(10000)
	txFee := int64(200)
	balance := int64(99800)

	_, err := CreateTransaction(privWif, txHash, destination, amount, txFee, balance)
	if err != nil {
		t.Error(err)
	}
}

func TestDecode(t *testing.T) {
	bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "btc", "test3"}
	skel, err := bcy.DecodeTX("0100000001283cf5ef49ce73a8630e08e3659e084effb92eef82ac0b647d7597c0a47707d3010000006b4830450221009dba1f31e237798be21afcf597bc880b6b4afc347421f2604b1ded4e05a8fef50220370a8c9a554a2fd16a8c401f5038f09ba6fd5bad9e3e5aec995555c7c6fefdb9012102da2c61d91c67d6615b3471b14d3d54bfbfd181537a2a231cdabcf62bd4fe3dfaffffffff01c85e0100000000001976a91443738e06bb02b07ad9c67e9480918d5df41fe35588ac00000000")
	if err != nil {
		fmt.Println(err)
		t.Error(err.Error())
	}
	fmt.Printf("%+v\n", skel)
}
