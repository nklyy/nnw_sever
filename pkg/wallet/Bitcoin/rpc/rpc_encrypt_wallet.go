package rpc

import "errors"

func EncryptWallet(password, walletName string) error {
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
		Method:  "encryptwallet",
		Params:  []string{password},
	}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return errors.New("could not sent transaction")
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}
