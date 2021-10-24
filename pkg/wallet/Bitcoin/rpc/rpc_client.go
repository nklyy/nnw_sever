package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func Client(body, res interface{}, walletInfo bool, walletName string) error {
	var serverAddr string

	if walletInfo {
		serverAddr = "http://138.68.92.142:8332/wallet/" + walletName // testnet/main net
	} else {
		serverAddr = "http://138.68.92.142:8332"
	}

	client := &http.Client{}

	jsonBody, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(jsonBody)
	req, err := http.NewRequest("POST", serverAddr, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("uuuset", "password123123")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, res)
	if err != nil {
		return err
	}

	if resp.Status != "200" {
		return errors.New(string(respBody))
	}

	return nil
}