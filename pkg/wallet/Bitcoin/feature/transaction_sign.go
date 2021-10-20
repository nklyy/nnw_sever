package feature

import (
	"errors"
)

func SignTx(tx string, privateKey string, unspentUtxos []*UnspentList) (string, error) {
	privateKeyArray := []string{privateKey}

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
		Params:  []interface{}{tx, privateKeyArray, unspentUtxos},
	}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return "", errors.New("could not sign transaction")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	if !msg.Result.Complete {
		return "", errors.New("tx hash not completed")
	}

	return msg.Result.Hex, nil
}
