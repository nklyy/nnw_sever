package rpc

import (
	"errors"
)

//type TxList struct {
//	Hash    string
//	TxIndex int64
//}

type Txs struct {
	Address  string      `json:"address"`
	Category string      `json:"category"`
	Amount   interface{} `json:"amount"`
	Txid     string      `json:"txid"`
}

func TransactionList(walletId string) ([]*UTXO, []*Txs, error) {
	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "listtransactions",
		Params:  []interface{}{"*", 1000000},
	}

	msg := struct {
		Result []struct {
			Address   string      `json:"address"`
			Category  string      `json:"category"`
			Amount    interface{} `json:"amount"`
			Vout      int64       `json:"vout"`
			Fee       interface{} `json:"fee"`
			Blockhash string      `json:"blockhash"`
			Txid      string      `json:"txid"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := Client(req, &msg, true, walletId)
	if err != nil {
		return nil, nil, errors.New("could not get transaction list")
	}

	if msg.Error.Message != "" {
		return nil, nil, errors.New(msg.Error.Message)
	}

	var utxos []*UTXO
	for idx := range msg.Result {
		utxos = append(utxos, &UTXO{
			TxId: msg.Result[idx].Txid,
			Vout: msg.Result[idx].Vout,
		})
	}

	var txs []*Txs
	for idx := range msg.Result {
		txs = append(txs, &Txs{
			Address:  msg.Result[idx].Address,
			Category: msg.Result[idx].Category,
			Amount:   msg.Result[idx].Amount,
			Txid:     msg.Result[idx].Txid,
		})
	}

	// Reverse arr, because we have old txs on first place, but we need new ones on first place
	for i, j := 0, len(txs)-1; i < j; i, j = i+1, j-1 {
		txs[i], txs[j] = txs[j], txs[i]
	}

	return utxos, txs, nil
}
