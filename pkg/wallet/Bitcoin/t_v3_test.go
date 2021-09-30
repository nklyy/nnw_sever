package Bitcoin

import (
	"fmt"
	"github.com/blockcypher/gobcy"
	"testing"
)

func TestRunTransactionV3(t *testing.T) {
	err := RunTransactionV3()
	if err != nil {
		t.Error(err)
	}
}

func TestDecode(t *testing.T) {
	bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "bcy", "test"}
	skel, err := bcy.DecodeTX("0100000001a5bb5619506f4e55b901d1550cb0926605dcb5c421c667e083faae33fb4844ea010000006a473044022025ae4060be871cd6474ce51f41c272857d6f01ef7a4f1a77e819d1224c96947f02207610b5ed202a260ef7163e18f0bf331758de1f26cdc72f05a94ba0a96e289a30012102c1e270dd70732f7e89f1ebe1c36e70240b57402530e741af0c12fcb7a34b7bd5ffffffff0110270000000000001976a91423241638a4d26167e906933ac8b03fe19a04d90e88ac00000000")
	if err != nil {
		fmt.Println(err)
		t.Error(err.Error())
	}
	fmt.Printf("%+v\n", skel)
}
