package rpc

import "errors"

func CreateWallet(walletName string) (string, error) {
	msg := struct {
		Result struct {
			Name string `json:"name"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "createwallet",
		Params:  []string{walletName},
	}

	err := RpcClient(req, &msg, false, "")
	if err != nil {
		return "", errors.New("could not sent transaction")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result.Name, nil
}
