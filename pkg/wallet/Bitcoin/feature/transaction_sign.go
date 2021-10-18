package feature

import (
	"errors"
)

func SignTx(tx string, privateKey string, utxos []*UTXO) (string, error) {
	privateKeyArray := []string{privateKey}

	var params []map[string]interface{}

	for idx := range utxos {
		paramMap := make(map[string]interface{})
		paramMap["txid"] = utxos[idx].Hash
		paramMap["vout"] = utxos[idx].TxIndex
		paramMap["scriptPubKey"] = utxos[idx].PKScript
		params = append(params, paramMap)
	}

	msg := struct {
		Result struct {
			Hex      string `json:"hex"`
			Complete bool   `json:"complete"`
		} `json:"result"`
		Error struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
		Id interface{} `json:"id"`
	}{}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "signrawtransactionwithkey",
		Params:  []interface{}{tx, privateKeyArray, params},
	}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return "", errors.New("could not get utxos")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	if !msg.Result.Complete {
		return "", errors.New("tx hash not completed")
	}

	return msg.Result.Hex, nil
}
