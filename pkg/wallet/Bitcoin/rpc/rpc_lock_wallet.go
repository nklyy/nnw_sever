package rpc

import "errors"

func LockWallet(walletName string) error {
	msg := struct {
		Result interface{} `json:"result"`
		Error  struct {
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
		Method:  "walletlock",
		Params:  []string{},
	}

	err := Client(req, &msg, true, walletName)
	if err != nil {
		return errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}
