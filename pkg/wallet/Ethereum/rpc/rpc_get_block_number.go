package rpc

import (
	"errors"
	"math/big"
	"nnw_s/pkg/helpers"
)

func GetBlockNumber() (*big.Int, error) {
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
		Method:  "eth_blockNumber",
		Params:  []string{},
		Id:      1,
	}

	err := Client(req, &msg)
	if err != nil {
		return nil, errors.New(msg.Error.Message)
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	res := helpers.ConvertHexToDecimal(msg.Result)

	return res, nil
}
