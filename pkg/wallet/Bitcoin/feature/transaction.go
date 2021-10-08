package feature

import "log"

func BuildTransaction() {
	//chainParams := &chaincfg.MainNetParams
	//chainParams := &chaincfg.TestNet3Params
	//
	//amountToSend := big.NewInt(50000) // amount to send in satoshis (0.01 btc)

	feeRate, err := GetCurrentFeeRate()
	log.Printf("current fee rate: %v", feeRate)
	if err != nil {
		log.Fatal(err)
	}

	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"

	log.Printf("from wallet public address: %s", fromWalletPublicAddress)

	_, err = ListUnspentTXOs(fromWalletPublicAddress)
	if err != nil {
		log.Fatal(err)
	}
}
