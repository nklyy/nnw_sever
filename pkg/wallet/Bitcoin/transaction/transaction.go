package transaction

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"log"
	"math/big"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

func BuildTransaction(fromWalletPublicAddress, destinationAddress, userWalletName, userWalletPassword string, amountToSend *big.Int) {
	//chainParams := &chaincfg.MainNetParams
	chainParams := &chaincfg.TestNet3Params

	// Get fee
	feeRate, err := rpc.GetCurrentFeeRate()
	log.Printf("%-18s %s\n", "current fee rate:", feeRate)
	if err != nil {
		log.Fatal(err)
	}

	// List unspent
	unspentTXOsList, err := rpc.ListUnspentTXOs(fromWalletPublicAddress, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate all unspent amount
	utxosAmount := big.NewInt(0)
	for idx := range unspentTXOsList {
		utxosAmount.Add(utxosAmount, unspentTXOsList[idx].Amount)
	}

	// Init transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// prepare transaction inputs
	sourceUtxosAmount := big.NewInt(0)
	var sourceUTXOs []*rpc.UnspentList
	for idx := range unspentTXOsList {
		hashStr := unspentTXOsList[idx].TxId
		sourceUtxosAmount.Add(sourceUtxosAmount, unspentTXOsList[idx].Amount)

		sourceUTXOHash, err := chainhash.NewHashFromStr(hashStr)
		if err != nil {
			log.Fatal(err)
		}

		if amountToSend.Int64() <= sourceUtxosAmount.Int64() {
			sourceUTXOIndex := uint32(unspentTXOsList[idx].Vout)
			sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
			sourceUTXOs = append(sourceUTXOs, &rpc.UnspentList{
				TxId:         unspentTXOsList[idx].TxId,
				Vout:         unspentTXOsList[idx].Vout,
				ScriptPubKey: unspentTXOsList[idx].PKScript,
			})
			sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)

			tx.AddTxIn(sourceTxIn)
			break
		}

		sourceUTXOIndex := uint32(unspentTXOsList[idx].Vout)
		sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
		sourceUTXOs = append(sourceUTXOs, &rpc.UnspentList{
			TxId:         unspentTXOsList[idx].TxId,
			Vout:         unspentTXOsList[idx].Vout,
			ScriptPubKey: unspentTXOsList[idx].PKScript,
		})
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)

		tx.AddTxIn(sourceTxIn)
	}

	// create the transaction outputs
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

	// calculate the change
	change := new(big.Int).Set(sourceUtxosAmount)
	change = new(big.Int).Sub(change, amountToSend)
	//change = new(big.Int).Sub(change, totalFee)
	if change.Cmp(big.NewInt(0)) == -1 {
		log.Fatal(err)
	}

	if change.Int64() != 0 {
		// our fee address
		//feeSendToAddress, err := btcutil.DecodeAddress(fromWalletPublicAddress, chainParams)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//feeSendToScript, err := txscript.PayToAddrScript(feeSendToAddress)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		////tx out to send change back to us
		//feeOutput := wire.NewTxOut(changeFee.Int64(), feeSendToScript)
		//tx.AddTxOut(feeOutput)

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
	}

	// calculate fees
	txByteSize := big.NewInt(int64(len(tx.TxIn)*180 + len(tx.TxOut)*34 + 10 + len(tx.TxIn)))
	totalFee := new(big.Int).Mul(feeRate, txByteSize)

	// Need add fee to spend amount and then compare
	if (amountToSend.Int64() - totalFee.Int64()) >= sourceUtxosAmount.Int64() {
		log.Fatal(errors.New("your balance too low for this transaction"))
	}

	log.Printf("%-18s %s\n", "total fee:", totalFee)

	// Change amount of source output transaction
	tx.TxOut[0].Value = amountToSend.Int64() - totalFee.Int64()

	// Transaction Hash
	notSignedTxBuf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(notSignedTxBuf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%-18s %s\n", "Not signed Tx:", hex.EncodeToString(notSignedTxBuf.Bytes()))

	// Prepare to sign tx
	// Unlock wallet
	err = rpc.UnLockWallet(userWalletPassword, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	// Get Private key
	privWif, err := rpc.GetAddressPrivateKey(fromWalletPublicAddress, userWalletName)
	if err != nil {
		log.Fatal(err)
	}

	decodedWif, err := btcutil.DecodeWIF(privWif)
	if err != nil {
		log.Fatal(err)
	}

	addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)
	if err != nil {
		log.Fatal(err)
	}

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	if err != nil {
		log.Fatal(err)
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		log.Fatal(err)
	}

	for i := range sourceUTXOs {
		sigScript, err := txscript.SignatureScript(tx, i, sourcePkScript, txscript.SigHashAll, decodedWif.PrivKey, true)
		if err != nil {
			log.Fatalf("could not generate pubSig; err: %v", err)
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	// Transaction Hash
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%-18s %s\n", "Redeem Tx:", hex.EncodeToString(buf.Bytes()))

	// Send Transaction
	sendHash, err := rpc.SendTx(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%-18s %s\n", "tx hash:", sendHash)
}
