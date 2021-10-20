package feature

import (
	"errors"
)

//type TxList struct {
//	Hash    string
//	TxIndex int64
//}

func TransactionList(walletName string) ([]*UTXO, error) {
	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "listtransactions",
		Params:  []interface{}{"*", 1},
	}

	msg := struct {
		Result []struct {
			Address   string      `json:"address"`
			Category  string      `json:"category"`
			Vout      int64       `json:"vout"`
			Fee       interface{} `json:"fee"`
			Blockhash string      `json:"blockhash"`
			Txid      string      `json:"txid"`
		} `json:"result"`
	}{}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return nil, errors.New("could not get transaction list")
	}

	var utxos []*UTXO
	for idx := range msg.Result {
		utxos = append(utxos, &UTXO{
			TxId: msg.Result[idx].Txid,
			Vout: msg.Result[idx].Vout,
		})
	}

	return utxos, nil
}
