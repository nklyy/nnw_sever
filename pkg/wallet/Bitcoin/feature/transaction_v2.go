package feature

import (
	"fmt"
	"log"
	"math/big"
)

func BuildTransactionV2(fromWalletPublicAddress, destinationAddress, userWalletName, userWalletPassword string, amountToSend *big.Int) {
	//chainParams := &chaincfg.TestNet3Params

	// Get smart fee
	feeRate, err := GetCurrentFeeRate()
	fmt.Printf("%-18s %v\n", "current fee rate:", feeRate)
	if err != nil {
		log.Fatal(err)
	}

	// Get list unspent tx
	utxos, err := ListUnspentTXOs(fromWalletPublicAddress, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	// Create Transaction
	createTxHash, unspentUtxosList, err := CreateTransaction(utxos, destinationAddress, amountToSend)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "TxHash:", createTxHash)
	fmt.Printf("%-18s %v\n", "UnspentParamList:", unspentUtxosList)

	// Fund for transaction
	fundTxHash, err := FundForTransaction(createTxHash, fromWalletPublicAddress, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "FundTxHash:", fundTxHash)

	// Unlock wallet
	err = UnLockWallet(userWalletPassword, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	// Get Private key
	privWif, err := GetAddressPrivateKey(fromWalletPublicAddress, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "Private key:", privWif)

	// Sign Transaction
	signTxHash, err := SignTx(fundTxHash, privWif, unspentUtxosList)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "SignTxHash:", signTxHash)

	// Send Transaction
	transactionHash, err := SendTx(signTxHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-18s %s\n", "TransactionHash:", transactionHash)
}
