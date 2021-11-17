package transaction

import (
	"errors"
	"math/big"
	"nnw_s/pkg/wallet/Bitcoin/rpc"
)

func SignAndSendTx(userWalletPassword, userWalletName, fromWalletPublicAddress, txHash string, amountToSend *big.Int) (string, error) {
	// List unspent
	unspentTXOsList, err := rpc.ListUnspentTXOs(fromWalletPublicAddress, userWalletName)
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

	err = rpc.UnLockWallet(userWalletPassword, userWalletName)
	if err != nil {
		return "", errors.New("Wrong password! ")
	}

	// Get Private key
	privWif, err := rpc.GetAddressPrivateKey(fromWalletPublicAddress, userWalletName)
	if err != nil {
		return "", err
	}

	// Sign Transaction
	signTxHash, err := rpc.SignTx(txHash, privWif, unspentTxs)
	if err != nil {
		return "", err
	}

	// Send Transaction
	transactionHash, err := rpc.SendTx(signTxHash)
	if err != nil {
		return "", err
	}

	return transactionHash, nil
}
