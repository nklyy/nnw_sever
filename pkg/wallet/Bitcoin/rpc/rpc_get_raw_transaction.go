package rpc

import (
	"errors"
	"fmt"
	"time"
)

// TODO: Add confirmation field, then check it on front-end and set label unconfirmed if confirmation less than 6 or confirm if confirmation bigger than 6.

type TxInfo struct {
	Txid string `json:"txid"`
	Vin  []struct {
		Txid string `json:"txid"`
		Vout int64  `json:"vout"`
	} `json:"vin"`
	Vout []struct {
		Value        float64 `json:"value"`
		N            int64   `json:"n"`
		ScriptPubKey struct {
			Addresses string `json:"address"`
		} `json:"scriptPubKey"`
	} `json:"vout"`
	Time          time.Time `json:"time"`
	Confirmations int64     `json:"confirmations"`
}

func GetRawTransaction(walletName, address, tx string) (*TxInfo, error) {
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
			Vin  []struct {
				Txid string `json:"txid"`
				Vout int64  `json:"vout"`
			} `json:"vin"`
			Vout []struct {
				Value        float64 `json:"value"`
				N            int64   `json:"n"`
				ScriptPubKey struct {
					Addresses string `json:"address"`
				} `json:"scriptPubKey"`
			} `json:"vout"`
			Time          int64 `json:"time"`
			Confirmations int64 `json:"confirmations"`
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

	return &TxInfo{
		Txid:          msg.Result.Txid,
		Vin:           msg.Result.Vin,
		Vout:          msg.Result.Vout,
		Time:          time.Unix(msg.Result.Time, 0),
		Confirmations: msg.Result.Confirmations,
	}, nil
}
