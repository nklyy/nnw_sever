package rpc

import (
	"errors"
	"github.com/btcsuite/btcutil"
	"math/big"
)

func GetBalance(walletId string) (*big.Int, error) {
	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "getbalance",
		Params:  []interface{}{"*", 6},
	}

	msg := struct {
		Result float64 `json:"result"`
		Error  struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := Client(req, &msg, true, walletId)
	if err != nil {
		return nil, errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	balance, _ := btcutil.NewAmount(msg.Result)

	return big.NewInt(int64(balance)), nil
}
