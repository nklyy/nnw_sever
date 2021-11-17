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
	"math/big"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

func CreateNotSignTx(fromWalletPublicAddress, destinationAddress, userWalletName string, amountToSend *big.Int) (string, *big.Int, error) {
	//chainParams := &chaincfg.MainNetParams
	chainParams := &chaincfg.TestNet3Params

	// Get fee
	feeRate, err := rpc.GetCurrentFeeRate()
	if err != nil {
		return "", nil, err
	}

	// List unspent
	unspentTXOsList, err := rpc.ListUnspentTXOs(fromWalletPublicAddress, userWalletName)
	if err != nil {
		return "", nil, err
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
			return "", nil, err
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
		return "", nil, err
	}

	destScript, err := txscript.PayToAddrScript(destAddress)
	if err != nil {
		return "", nil, err
	}

	// tx out to send btc to user
	destOutput := wire.NewTxOut(amountToSend.Int64(), destScript)
	tx.AddTxOut(destOutput)

	// calculate the change
	change := new(big.Int).Set(sourceUtxosAmount)
	change = new(big.Int).Sub(change, amountToSend)
	//change = new(big.Int).Sub(change, totalFee)
	//if change.Cmp(big.NewInt(0)) == -1 {
	//	return "", nil, err
	//}

	if change.Int64() != 0 {
		// our change address
		changeSendToAddress, err := btcutil.DecodeAddress(fromWalletPublicAddress, chainParams)
		if err != nil {
			return "", nil, err
		}

		changeSendToScript, err := txscript.PayToAddrScript(changeSendToAddress)
		if err != nil {
			return "", nil, err
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
		return "", nil, errors.New("your balance too low for this transaction")
	}

	// Change amount of source output transaction
	tx.TxOut[0].Value = amountToSend.Int64() - totalFee.Int64()

	// Transaction Hash
	notSignedTxBuf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(notSignedTxBuf)
	if err != nil {
		return "", nil, err
	}

	return hex.EncodeToString(notSignedTxBuf.Bytes()), totalFee, nil
}
