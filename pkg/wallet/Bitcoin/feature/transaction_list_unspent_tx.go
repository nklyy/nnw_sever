package feature

import (
	"errors"
	"fmt"
	"log"
	"math/big"
)

// UTXO ...
type UTXO struct {
	Hash      string
	TxIndex   int
	Amount    *big.Int
	Spendable bool
	PKScript  []byte
}

func ListUnspentTXOs(address string) ([]*UTXO, error) {
	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "listunspent",
		Params:  []string{address},
	}

	fmt.Println(req)

	msg := struct {
		Result []struct {
			Txid          string   `json:"txid"`
			Vout          int64    `json:"vout"`
			Address       string   `json:"address"`
			Label         string   `json:"label"`
			ScriptPubKey  string   `json:"script_pub_key"`
			Amount        *big.Int `json:"amount"`
			Confirmations int64    `json:"confirmations"`
			RedeemScript  string   `json:"redeem_script"`
			WitnessScript string   `json:"witness_script"`
			Spendable     bool     `json:"spendable"`
			Solvable      bool     `json:"solvable"`
			Reused        bool     `json:"reused"`
			Desc          string   `json:"desc"`
			Safe          bool     `json:"safe"`
		} `json:"result"`
	}{}

	err := RpcClient(req, &msg)
	if err != nil {
		log.Printf("could not get utxos")
		return nil, errors.New("could not get utxos")
	}

	var utxos []*UTXO
	for idx := range msg.Result {
		utxos = append(utxos, &UTXO{
			//Hash:      msg.Result[idx].TXHash,
			//TxIndex:   int(msg.Result[idx].TXPosition),
			Amount:    msg.Result[idx].Amount,
			Spendable: msg.Result[idx].Spendable,
		})
	}

	return utxos, nil
}
