package rpc

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"strings"
)

func FundForTransaction(createTxHash, changeAddress, walletName string) (string, error) {
	subtractFeeFromOutputs := []int64{0}

	params := map[string]interface{}{
		"changeAddress":          changeAddress,
		"subtractFeeFromOutputs": subtractFeeFromOutputs,
	}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "fundrawtransaction",
		Params:  []interface{}{createTxHash, params},
	}

	msg := struct {
		Result struct {
			Hex string  `json:"hex"`
			Fee float64 `json:"fee"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return "", errors.New("could not fund for transaction")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	feeAmount, err := btcutil.NewAmount(msg.Result.Fee)
	if err != nil {
		return "", err
	}

	fmt.Printf("%-18s %v\n", "Fee for transaction in BTC:", feeAmount)
	fmt.Printf("%-18s %v\n", "Fee for transaction in Satoshi:", feeAmount.Format(btcutil.AmountSatoshi))
	fmt.Println(strings.Repeat("-", 106))

	return msg.Result.Hex, nil
}
