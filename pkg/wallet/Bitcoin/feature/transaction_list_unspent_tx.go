package feature

import (
	"errors"
	"fmt"
)

// UTXO ...
type UTXO struct {
	Hash      string
	TxIndex   int
	Amount    interface{}
	Spendable bool
	PKScript  []byte
}

func ListUnspentTXOs(address string) ([]*UTXO, error) {
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
			Txid          string      `json:"txid"`
			Vout          int64       `json:"vout"`
			Address       string      `json:"address"`
			Label         string      `json:"label"`
			ScriptPubKey  string      `json:"scriptPubKey"`
			Amount        interface{} `json:"amount"`
			Confirmations int64       `json:"confirmations"`
			RedeemScript  string      `json:"redeem_script"`
			WitnessScript string      `json:"witness_script"`
			Spendable     bool        `json:"spendable"`
			Solvable      bool        `json:"solvable"`
			Reused        bool        `json:"reused"`
			Desc          string      `json:"desc"`
			Safe          bool        `json:"safe"`
		} `json:"result"`
	}{}

	err := RpcClient(req, &msg, true, "first")
	if err != nil {
		return nil, errors.New("could not get utxos")
	}

	var utxos []*UTXO
	for idx := range msg.Result {
		utxos = append(utxos, &UTXO{
			Hash:      msg.Result[idx].Txid,
			TxIndex:   int(msg.Result[idx].Vout),
			Amount:    msg.Result[idx].Amount,
			Spendable: msg.Result[idx].Spendable,
		})
	}

	fmt.Println(utxos)

	return utxos, nil
}
