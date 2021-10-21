package rpc

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type UnspentList struct {
	TxId         string `json:"txid"`
	Vout         int64  `json:"vout"`
	ScriptPubKey string `json:"scriptPubKey"`
}

func CreateTransaction(utxos []*UTXO, addressTo string, spendAmount *big.Int) (string, []*UnspentList, error) {
	//var unspentParams []map[string]interface{}
	var unspentParams []*UnspentList

	utxosAmount := big.NewInt(0)

	// Convert spendAmount to BTC amount
	//fmt.Println(float64(spendAmount.Int64()) / 1.0e8)

	// Add all unspent amount
	for idx := range utxos {
		utxosAmount.Add(utxosAmount, utxos[idx].Amount)
	}

	// Need add fee to spend amount and then compare
	if spendAmount.Int64() >= utxosAmount.Int64() {
		return "", nil, errors.New("your balance too low for this transaction")
	}

	sourceUtxosAmount := big.NewInt(0)
	for idx := range utxos {
		sourceUtxosAmount.Add(sourceUtxosAmount, utxos[idx].Amount)

		if spendAmount.Int64() > sourceUtxosAmount.Int64() {
			//paramMap := make(map[string]interface{})
			//paramMap["txid"] = utxos[idx].TxId
			//paramMap["vout"] = utxos[idx].Vout
			//paramMap["scriptPubKey"] = utxos[idx].PKScript
			//unspentParams = append(unspentParams, paramMap)

			unspentParams = append(unspentParams, &UnspentList{
				TxId:         utxos[idx].TxId,
				Vout:         utxos[idx].Vout,
				ScriptPubKey: utxos[idx].PKScript,
			})
		} else {
			unspentParams = append(unspentParams, &UnspentList{
				TxId:         utxos[idx].TxId,
				Vout:         utxos[idx].Vout,
				ScriptPubKey: utxos[idx].PKScript,
			})
			break
		}
	}

	addressesParams := []interface{}{
		map[string]interface{}{addressTo: float64(spendAmount.Int64()) / 1.0e8},
		//map[string]interface{}{"mrvZjXUNupoEpQf6KsgiVgzLz7DUca6Kfv": 0.00001000},
		//map[string]interface{}{addressFrom: utxosAmount.Sub(utxosAmount, spendAmount)},
	}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "createrawtransaction",
		Params:  []interface{}{unspentParams, addressesParams},
	}

	msg := struct {
		Result string `json:"result"`
		Error  struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return "", nil, errors.New("could not create transaction")
	}

	if msg.Error.Message != "" {
		return "", nil, errors.New(msg.Error.Message)
	}

	fmt.Printf("%-18s %s\n", "Balance:", utxosAmount)
	fmt.Printf("%-18s %s\n", "Spend amount:", spendAmount)
	fmt.Printf("%-18s %s\n", "Remainder: ", utxosAmount.Sub(utxosAmount, spendAmount))
	fmt.Println(strings.Repeat("-", 106))

	return msg.Result, unspentParams, nil
}
