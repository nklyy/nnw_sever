package feature

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"log"
	"math/big"
)

func BuildTransaction() {
	//chainParams := &chaincfg.MainNetParams
	chainParams := &chaincfg.TestNet3Params

	amountToSend := big.NewInt(10000) // amount to send in satoshis (0.01 btc)

	feeRate, err := GetCurrentFeeRate()
	log.Printf("current fee rate: %v", feeRate)
	if err != nil {
		log.Fatal(err)
	}

	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"

	log.Printf("from wallet public address: %s", fromWalletPublicAddress)

	unspentTXOs, err := ListUnspentTXOs(fromWalletPublicAddress)
	if err != nil {
		log.Fatal(err)
	}

	balance := big.NewInt(110700)
	unspentTXOs, UTXOsAmount, err := marshalUTXOs(unspentTXOs, amountToSend, feeRate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(unspentTXOs[0])
	log.Println("unspent UTXOs", unspentTXOs, UTXOsAmount)

	tx := wire.NewMsgTx(wire.TxVersion)

	var sourceUTXOs []*UTXO
	// prepare tx ins
	for idx := range unspentTXOs {
		unspentTXOs[idx].Amount = balance
		hashStr := unspentTXOs[idx].Hash

		sourceUTXOHash, err := chainhash.NewHashFromStr(hashStr)
		if err != nil {
			log.Fatal(err)
		}

		sourceUTXOIndex := uint32(unspentTXOs[idx].TxIndex)
		sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
		sourceUTXOs = append(sourceUTXOs, unspentTXOs[idx])
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)

		tx.AddTxIn(sourceTxIn)
	}

	destinationAddress := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"

	// create the tx outs
	destAddress, err := btcutil.DecodeAddress(destinationAddress, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	destScript, err := txscript.PayToAddrScript(destAddress)
	if err != nil {
		log.Fatal(err)
	}

	// tx out to send btc to user
	destOutput := wire.NewTxOut(amountToSend.Int64(), destScript)
	tx.AddTxOut(destOutput)

	// calculate fees
	txByteSize := big.NewInt(int64(len(tx.TxIn)*180 + len(tx.TxOut)*34 + 10 + len(tx.TxIn)))
	fmt.Println(txByteSize.Int64(), "ASDASDASDASDAS")
	totalFee := new(big.Int).Mul(feeRate, txByteSize)
	log.Printf("total fee: %s", totalFee)

	//avByte := big.NewInt(225)
	//totalFee := new(big.Int).Mul(feeRate, avByte)
	//fmt.Println("TOTALFEE", totalFee)

	// calculate the change
	change := new(big.Int).Set(balance)
	change = new(big.Int).Sub(change, amountToSend)
	change = new(big.Int).Sub(change, totalFee)
	if change.Cmp(big.NewInt(0)) == -1 {
		log.Fatal(err)
	}

	// our change address
	changeSendToAddress, err := btcutil.DecodeAddress(fromWalletPublicAddress, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	changeSendToScript, err := txscript.PayToAddrScript(changeSendToAddress)
	if err != nil {
		log.Fatal(err)
	}

	//tx out to send change back to us
	changeOutput := wire.NewTxOut(change.Int64(), changeSendToScript)
	tx.AddTxOut(changeOutput)

	privWif := "cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD"

	//decodedWif, err := btcutil.DecodeWIF(privWif)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Printf("Source Address: %s\n", sourceAddress) // Source Address: mgjHgKi1g6qLFBM1gQwuMjjVBGMJdrs9pP

	//sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for i := range sourceUTXOs {
	//	sigScript, err := txscript.SignatureScript(tx, i, sourcePkScript, txscript.SigHashAll, decodedWif.PrivKey, true)
	//	if err != nil {
	//		log.Fatalf("could not generate pubSig; err: %v", err)
	//	}
	//	tx.TxIn[i].SignatureScript = sigScript
	//}

	// Transaction Hash
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes()))

	t := hex.EncodeToString(buf.Bytes())

	signTx, err := SignTx(t, privWif, sourceUTXOs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SignedTX:", signTx)

	sendHash, err := SendTx(signTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx hash: %s\n", sendHash) // 1d8f70dfc8b90bff672ee663a7cc811c4e88e98c6895dc93aa9f73202bb7809b
}
