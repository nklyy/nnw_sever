package rpc

import (
	"errors"
	"strconv"
	"strings"
)

func GetBalance(address string) (*int64, error) {
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
		Method:  "eth_getBalance",
		Params:  []string{address, "latest"},
		Id:      1,
	}

	err := Client(req, &msg)
	if err != nil {
		return nil, errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	numberStr := strings.Replace(msg.Result, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	balance, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
