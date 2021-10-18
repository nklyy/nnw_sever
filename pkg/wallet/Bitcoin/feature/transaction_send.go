package feature

import "errors"

func SendTx(signedTX string) (string, error) {
	msg := struct {
		Result struct {
			Hex string `json:"hex"`
		} `json:"result"`
		Error struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
		Id interface{} `json:"id"`
	}{}

	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "sendrawtransaction",
		Params:  []string{signedTX},
	}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return "", errors.New("could not get utxos")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result.Hex, nil
}
