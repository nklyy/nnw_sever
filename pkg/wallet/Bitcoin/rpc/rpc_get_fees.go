package rpc

import (
	"errors"
	"log"
	"math/big"
)

// GetCurrentFee gets the current fee in bitcoin
func GetCurrentFee() (*float64, error) {
	req := struct {
		JsonRPC string `json:"json_rpc"`
		Method  string `json:"method"`
		Params  []int  `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "estimatesmartfee",
		Params:  []int{2},
	}

	msg := struct {
		Result struct {
			Feerate float64 `json:"feerate"`
			Blocks  int64   `json:"blocks"`
		}
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := RpcClient(req, &msg, false, "first")
	if err != nil {
		return nil, err
	}

	if msg.Result.Feerate == -1.0 || msg.Result.Feerate == 0 {
		log.Printf("expected result > 0; received: %f", msg.Result.Feerate)
	}

	if msg.Error.Message != "" {
		return nil, errors.New(msg.Error.Message)
	}

	var fee float64
	fee = msg.Result.Feerate
	// sanity check
	if fee > 0.05 {
		fee = 0.1
	} else if fee < 0 {
		fee = 0
	}
	//fmt.Printf("fee: %f\n", fee)

	if fee == 0 {
		log.Print("could not get fees")
		return &fee, errors.New("could not get fees")
	}

	return &fee, nil
}

// GetCurrentFeeRate gets the current fee in satoshis per kb
func GetCurrentFeeRate() (*big.Int, error) {
	fee, err := GetCurrentFee()
	if err != nil {
		return nil, err
	}

	// convert to satoshis to bytes
	// feeRate := big.NewInt(int64(msg.Result * 1.0E8))
	// convert to satoshis to kb
	feeRate := big.NewInt(int64(*fee * 1.0e5))

	//fmt.Printf("fee rate: %s\n", feeRate)

	return feeRate, nil
}
