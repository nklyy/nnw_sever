package rpc

import (
	"errors"
)

func LoadWallet(walletID string) (string, error) {
	msg := struct {
		Result struct {
			Name    string `json:"name"`
			Warning string `json:"warning"`
		} `json:"result"`
		Error struct {
			Code    int64  `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}

	req := struct {
		JsonRPC string   `json:"json_rpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "loadwallet",
		Params:  []string{walletID},
	}

	err := Client(req, &msg, true, walletID)

	if msg.Error.Code == -4 || msg.Error.Code == -35 {
		return "", nil
	}

	if err != nil {
		return "", errors.New("could not sent transaction")
	}

	if msg.Error.Message != "" {
		return "", errors.New(msg.Error.Message)
	}

	return msg.Result.Warning, nil
}
