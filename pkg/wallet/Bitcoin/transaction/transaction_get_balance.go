package transaction

import (
	"errors"
	"github.com/btcsuite/btcutil"
	"math/big"
)

func GetBalance(walletName string) (*big.Int, error) {
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

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return nil, errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	balance, _ := btcutil.NewAmount(msg.Result)

	return big.NewInt(int64(balance)), nil
}
