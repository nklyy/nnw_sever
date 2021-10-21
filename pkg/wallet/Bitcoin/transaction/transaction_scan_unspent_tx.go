package transaction

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"math/big"
	"strconv"
)

func ScanUnspentTXOs(address string) ([]*UTXO, error) {
	addressArray := []string{"addr(" + address + ")"}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "scantxoutset",
		Params:  []interface{}{"start", addressArray},
	}

	msg := struct {
		Result struct {
			Success   bool   `json:"success"`
			Txouts    int64  `json:"txouts"`
			Height    int64  `json:"height"`
			Bestblock string `json:"bestblock"`
			Unspents  []struct {
				Txid         string
				Vout         int64
				ScriptPubkey string
				Desc         string
				Amount       float64
				Height       int64
			} `json:"unspents"`
			TotalAmount interface{} `json:"total_amount"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return nil, errors.New("could not scan unspent transactions")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	var utxos []*UTXO
	for idx := range msg.Result.Unspents {
		btcAmount, _ := btcutil.NewAmount(msg.Result.Unspents[idx].Amount)

		utxos = append(utxos, &UTXO{
			TxId:     msg.Result.Unspents[idx].Txid,
			Vout:     msg.Result.Unspents[idx].Vout,
			Amount:   big.NewInt(int64(btcAmount)),
			PKScript: msg.Result.Unspents[idx].ScriptPubkey,
			//Spendable: msg.Result[idx].Spendable,
		})
	}

	totalStr := fmt.Sprintf("%v", msg.Result.TotalAmount)
	float, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return nil, err
	}
	totalAmount, _ := btcutil.NewAmount(float)
	fmt.Println("TOTAL", int64(totalAmount))

	return utxos, nil
}
