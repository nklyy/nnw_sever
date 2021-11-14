package rpc

import (
	"errors"
	"fmt"
)

//type TxList struct {
//	Hash    string
//	TxIndex int64
//}

type TxSimple struct {
	ToAddress   []string `json:"to_address"`
	FromAddress []string `json:"from_address"`
	Amount      float64  `json:"amount"`
	TxId        string   `json:"tx_id"`
}

func GetRawTransaction(walletName, address, tx string) ([]*Txs, error) {
	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "getrawtransaction",
		Params:  []interface{}{tx, true},
	}

	msg := struct {
		Result struct {
			Txid string `json:"txid"`
			Vout []struct {
				Value        float64 `json:"value"`
				N            int64   `json:"n"`
				ScriptPubKey struct {
					Addresses []string `json:"addresses"`
				} `json:"scriptPubKey"`
			} `json:"vout"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := Client(req, &msg, true, walletName)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("could not get transaction list")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	fmt.Println(msg)

	return nil, nil
}
