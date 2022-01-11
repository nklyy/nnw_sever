package transaction

import (
	"errors"
	"math/big"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

func SignAndSendTx(userWalletPassword, userWalletId, fromWalletPublicAddress, txHash string, amountToSend *big.Int) (string, error) {
	// List unspent
	unspentTXOsList, err := rpc.ListUnspentTXOs(fromWalletPublicAddress, userWalletId)
	if err != nil {
		return "", err
	}

	// Calculate all unspent amount
	utxosAmount := big.NewInt(0)
	for idx := range unspentTXOsList {
		utxosAmount.Add(utxosAmount, unspentTXOsList[idx].Amount)
	}

	// prepare transaction inputs
	sourceUtxosAmount := big.NewInt(0)
	var unspentTxs []*rpc.UnspentList
	for idx := range unspentTXOsList {
		sourceUtxosAmount.Add(sourceUtxosAmount, unspentTXOsList[idx].Amount)

		if amountToSend.Int64() <= sourceUtxosAmount.Int64() {
			unspentTxs = append(unspentTxs, &rpc.UnspentList{
				TxId:         unspentTXOsList[idx].TxId,
				Vout:         unspentTXOsList[idx].Vout,
				ScriptPubKey: unspentTXOsList[idx].PKScript,
			})
			break
		}

		unspentTxs = append(unspentTxs, &rpc.UnspentList{
			TxId:         unspentTXOsList[idx].TxId,
			Vout:         unspentTXOsList[idx].Vout,
			ScriptPubKey: unspentTXOsList[idx].PKScript,
		})
	}

	err = rpc.UnLockWallet(userWalletPassword, userWalletId)
	if err != nil {
		return "", errors.New("Wrong password! ")
	}

	// Get Private key
	privateKey, err := rpc.GetAddressPrivateKey(fromWalletPublicAddress, userWalletId)
	if err != nil {
		return "", err
	}

	// Sign Transaction
	signTxHash, err := rpc.SignTx(txHash, privateKey, unspentTxs)
	if err != nil {
		return "", err
	}

	// Send Transaction
	transactionHash, err := rpc.SendTx(signTxHash)
	if err != nil {
		return "", err
	}

	err = rpc.LockWallet(userWalletId)
	if err != nil {
		return "", err
	}

	return transactionHash, nil
}
