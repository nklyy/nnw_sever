package rpc

import "errors"

func GetTransactionCount(address string) error {
	msg := struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int64  `json:"id"`
		Result  string `json:"result"`
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
		Method:  "eth_getTransactionCount",
		Params:  []string{address, "latest"},
		Id:      1,
	}

	err := Client(req, &msg)
	if err != nil {
		return errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}
