package rpc

import (
	"errors"
)

func AddressInfo(address, walletId string) (string, error) {
	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "getaddressinfo",
		Params:  []string{address},
	}

	msg := struct {
		Result struct {
			Address      string `json:"address"`
			ScriptPubKey string `json:"scriptPubKey"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}{}

	err := Client(req, &msg, true, walletId)
	if err != nil {
		return "", errors.New("could not get address info")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result.ScriptPubKey, nil
}
