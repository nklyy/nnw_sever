package rpc

import (
	"errors"
)

func UnLockWallet(password, walletId string) error {
	msg := struct {
		Result interface{} `json:"result"`
		Error  struct {
			Code    int64  `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "walletpassphrase",
		Params:  []interface{}{password, 60},
	}

	err := Client(req, &msg, true, walletId)
	if err != nil {
		return errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}
