package feature

import "errors"

func AddressInfo(address, walletName string) (string, error) {
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
	}{}

	err := RpcClient(req, &msg, true, walletName)
	if err != nil {
		return "", errors.New("could not get address info")
	}

	return msg.Result.ScriptPubKey, nil
}
