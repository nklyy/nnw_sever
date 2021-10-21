package rpc

import (
	"errors"
)

func GetAddressPrivateKey(address, walletName string) (string, error) {
	msg := struct {
		Result string `json:"result"`
		Error  struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "dumpprivkey",
		Params:  []string{address},
	}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return "", errors.New("could not sent transaction")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result, nil
}
