package rpc

import "errors"

type SendTx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func SendTransaction(from, to, value string) (string, error) {
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
		Params  []SendTx `json:"params"`
		Id      int64    `json:"id"`
	}{
		JsonRPC: "2.0",
		Method:  "eth_sendTransaction",
		Params: []SendTx{
			{
				From:  from,
				To:    to,
				Value: value,
			},
		},
		Id: 1,
	}

	err := Client(req, &msg)
	if err != nil {
		return "", errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result, nil
}
