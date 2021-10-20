package feature

import "errors"

func UnLockWallet(password, walletName string) error {
	msg := struct {
		Result interface{} `json:"result"`
	}{}

	req := struct {
		JsonRPC string        `json:"json_rpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JsonRPC: "2.0",
		Method:  "walletpassphrase",
		Params:  []interface{}{password, 60},
	}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return errors.New("could not sent transaction")
	}

	return nil
}
