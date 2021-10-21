package not_working

import (
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	privWif := "cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD"
	txHash := "f9a863929bfa7bc777e19795453f78a3867090118d93fbbc431d0089f88a27a0"
	destination := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
	amount := int64(10000)
	txFee := int64(300)
	balance := int64(20444)

	_, err := CreateTransaction(privWif, txHash, destination, amount, txFee, balance)
	if err != nil {
		t.Error(err)
	}
}

//
//func TestDecode(t *testing.T) {
//	bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "btc", "test3"}
//	skel, err := bcy.DecodeTX("01000000010615165ff7a6f08ec0c37e98faa32e9920cb40fe2510faec0013129293fbf160010000006a47304402200cfe88a85dbd635adb53d8e2b455c84b9f8ea9f3425dd6adcb87f5a0eda7d97e022062e0a16fbd13ac4931100c8ef0e7438d3192fa93e85f995201eef1d52e78933b012102da2c61d91c67d6615b3471b14d3d54bfbfd181537a2a231cdabcf62bd4fe3dfaffffffff0210270000000000001976a91443738e06bb02b07ad9c67e9480918d5df41fe35588acc00c0100000000001976a914690cd6356789d30b99063632e0651a8d0c206c7f88ace8030000")
//	if err != nil {
//		fmt.Println(err)
//		t.Error(err.Error())
//	}
//	fmt.Printf("%+v\n", skel)
//}
