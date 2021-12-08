package rpc

import "errors"

func LockWallet(address string) (bool, error) {
	msg := struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int64  `json:"id"`
		Result  bool   `json:"result"`
		Error   struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}{}

	req := struct {
		JsonRPC string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		Id      int64    `json:"id"`
	}{
		JsonRPC: "2.0",
		Method:  "personal_lockAccount",
		Params:  []string{address},
		Id:      1,
	}

	err := Client(req, &msg)
	if err != nil {
		return false, errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return false, errors.New(msg.Error.Message)
	}

	return msg.Result, nil
}
