package rpc

import (
	"errors"
	"github.com/btcsuite/btcutil"
	"math/big"
)

// UTXO ...
type UTXO struct {
	TxId   string
	Vout   int64
	Amount *big.Int
	//Spendable bool
	PKScript string
}

func ListUnspentTXOs(address string, walletName string) ([]*UTXO, error) {
	addressArray := []string{address}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "listunspent",
		Params:  []interface{}{1, 9999999, addressArray},
	}

	msg := struct {
		Result []struct {
			Txid          string  `json:"txid"`
			Vout          int64   `json:"vout"`
			Address       string  `json:"address"`
			Label         string  `json:"label"`
			ScriptPubKey  string  `json:"scriptPubKey"`
			Amount        float64 `json:"amount"`
			Confirmations int64   `json:"confirmations"`
			RedeemScript  string  `json:"redeem_script"`
			WitnessScript string  `json:"witness_script"`
			Spendable     bool    `json:"spendable"`
			Solvable      bool    `json:"solvable"`
			Reused        bool    `json:"reused"`
			Desc          string  `json:"desc"`
			Safe          bool    `json:"safe"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return nil, errors.New("could not get utxos")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	var utxos []*UTXO
	for idx := range msg.Result {
		btcAmount, _ := btcutil.NewAmount(msg.Result[idx].Amount)

		utxos = append(utxos, &UTXO{
			TxId:     msg.Result[idx].Txid,
			Vout:     msg.Result[idx].Vout,
			Amount:   big.NewInt(int64(btcAmount)),
			PKScript: msg.Result[idx].ScriptPubKey,
			//Spendable: msg.Result[idx].Spendable,
		})
	}

	return utxos, nil
}
