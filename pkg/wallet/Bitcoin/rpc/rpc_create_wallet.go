package rpc

import (
	"errors"
)

func CreateWallet(walletId string) (string, error) {
	msg := struct {
		Result struct {
			Name string `json:"name"`
		} `json:"result"`
		Error struct {
			Code    int64  `json:"code"`
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
		Params:  []string{walletId},
	}

	err := Client(req, &msg, false, "")
	if err != nil {
		return "", errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result.Name, nil
}
