package rpc

import (
	"errors"
	"math/big"
	"nnw_s/pkg/helpers"
)

func GetBalance(address string) (*big.Int, error) {
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
		// TODO use uuid in Id field and then check it in result
		Id: 1,
	}

	err := Client(req, &msg)
	if err != nil {
		return nil, errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	res := helpers.ConvertHexToDecimal(msg.Result)

	//balance, err := strconv.ParseInt(numberStr, 16, 64)
	//if err != nil {
	//	return nil, err
	//}

	return res, nil
}
