package feature

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"math/big"
	"strconv"
)

// UTXO ...
type UTXO struct {
	Hash      string
	TxIndex   int64
	Amount    *big.Int
	Spendable bool
	PKScript  string
}

func ListUnspentTXOs(address string) ([]*UTXO, error) {
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

	//msg := struct {
	//	Result []struct {
	//		Txid          string      `json:"txid"`
	//		Vout          int64       `json:"vout"`
	//		Address       string      `json:"address"`
	//		Label         string      `json:"label"`
	//		ScriptPubKey  string      `json:"scriptPubKey"`
	//		Amount        interface{} `json:"amount"`
	//		Confirmations int64       `json:"confirmations"`
	//		RedeemScript  string      `json:"redeem_script"`
	//		WitnessScript string      `json:"witness_script"`
	//		Spendable     bool        `json:"spendable"`
	//		Solvable      bool        `json:"solvable"`
	//		Reused        bool        `json:"reused"`
	//		Desc          string      `json:"desc"`
	//		Safe          bool        `json:"safe"`
	//	} `json:"result"`
	//}{}

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
	}{}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return nil, errors.New("could not get utxos")
	}

	var utxos []*UTXO
	for idx := range msg.Result.Unspents {
		btcAmount, _ := btcutil.NewAmount(msg.Result.Unspents[idx].Amount)

		utxos = append(utxos, &UTXO{
			Hash:     msg.Result.Unspents[idx].Txid,
			TxIndex:  msg.Result.Unspents[idx].Vout,
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
