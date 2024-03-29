package not_working

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func CreateTransaction(privWif string, txHash string, destination string, amount int64, txFee int64, balance int64) (string, error) {
	sourceUTXOIndex := uint32(0)
	chainParams := &chaincfg.TestNet3Params

	decodedWif, err := btcutil.DecodeWIF(privWif)
	if err != nil {
		return "", err
	}

	fmt.Printf("%-18s %v\n", "Decoded WIF: ", decodedWif) // Decoded WIF: cTujmQgVdGYzmZEfhq5gVDpd2EAHF1sZahPmkDnRHmPDEVRYz6eo

	addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)
	if err != nil {
		return "", err
	}

	sourceUTXOHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", err
	}

	fmt.Printf("%-18s %s\n", "UTXO hash: ", sourceUTXOHash) // utxo hash: c6950f355835c361dce2e9d6eb511cf56972b67cb34dad5d1fd9f9bc796711a5

	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(destination, chainParams)
	if err != nil {
		return "", err
	}

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	if err != nil {
		return "", err
	}

	fmt.Printf("%-18s %s\n", "Source Address: ", sourceAddress) // Source Address: mqJ8FALtYnxvLgwTUWQ2shNkdiLuU7tkPR

	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return "", err
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return "", err
	}

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	//redeemTx.LockTime = 2097025
	redeemTx.AddTxIn(sourceTxIn)

	redeemTxOut := wire.NewTxOut(amount-txFee, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	//
	//redeemTxOut = wire.NewTxOut(balance-amount-txFee, sourcePkScript)
	//redeemTx.AddTxOut(redeemTxOut)

	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTxOut.PkScript, txscript.SigHashAll, decodedWif.PrivKey, true)
	if err != nil {
		return "", err
	}

	redeemTx.TxIn[0].SignatureScript = sigScript
	fmt.Printf("%-18s %v\n", "Signature Script: ", hex.EncodeToString(sigScript)) // Signature Script: 473...b67

	// validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTxOut.PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return "", err
	}

	if err := vm.Execute(); err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(make([]byte, 0, redeemTx.SerializeSize()))
	err = redeemTx.Serialize(buf)
	if err != nil {
		return "", err
	}

	fmt.Printf("%-18s %v\n", "Redeem Tx: ", hex.EncodeToString(buf.Bytes())) // redeem Tx: 01000000011efc...5bb88ac00000000

	//execute transaction
	//curl -v --digest -X POST 127.0.0.1:8332 -d "{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"getblockchaininfo\"}"-H 'Content-Type:application/json'
	//connConfig := &rpcclient.ConnConfig{
	//	HTTPPostMode: true,
	//	DisableTLS:   true,
	//	Host:         "127.0.0.1:8332",
	//	User:         "uuuset",
	//	Pass:         "password123123",
	//}
	//bitcoinClient, err := rpcclient.New(connConfig, nil)
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Println("CONNECTED")
	//asd, err := bitcoinClient.GetBlockChainInfo()
	//if err != nil {
	//	fmt.Println("ERRRRRRR")
	//	return "", err
	//}
	//fmt.Println(asd)

	//a := btcjson.TransactionInput{
	//	Txid: "92ddc3ab2b148e1c8a180070fc8d4c49496c88eea565a78d54907a3848fcf786",
	//	Vout: 0,
	//}
	//
	//var qwe []btcjson.TransactionInput
	//qwe = append(qwe, a)
	//
	//ammmm, _ := btcutil.NewAmount(0.000915)
	//ammmm1, _ := btcutil.NewAmount(0.00008)
	//
	//out := map[btcutil.Address]btcutil.Amount{destinationAddress: ammmm, sourceAddress: ammmm1}
	//
	//transaction, err := bitcoinClient.CreateRawTransaction(qwe, out, nil)
	//if err != nil {
	//	return "", err
	//}
	//
	//buff := bytes.NewBuffer(make([]byte, 0, transaction.SerializeSize()))
	//transaction.Serialize(buff)
	//
	//fmt.Println("TRRRRRRRRRRRRRRR", hex.EncodeToString(buff.Bytes()))

	//decode, err := bitcoinClient.DecodeRawTransaction(buf.Bytes())
	//if err != nil {
	//	return "", err
	//}
	//fmt.Println("DECODE", decode)
	//
	//bal, err := bitcoinClient.GetBalance("mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u")
	//if err != nil {
	//	return "", err
	//}
	//println("hash:", bal)

	//hash, err := bitcoinClient.SendRawTransaction(redeemTx, false)
	//if err != nil {
	//	return "", err
	//}
	//println("hash:", hash.String())

	// Push Transaction
	//bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "btc", "test3"}
	//skel, err := bcy.PushTX(hex.EncodeToString(buf.Bytes()))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%+v\n", skel)

	// Make POST request
	//url := "https://chain.so/api/v2/send_tx/BTCTEST"
	//url := "https://mempool.space/testnet/api/tx"
	//
	//postBody, _ := json.Marshal(map[string]string{
	//	//"network": "BTCTEST",
	//	//"tx_hex":  hex.EncodeToString(buf.Bytes()),
	//	"txHex": hex.EncodeToString(buf.Bytes()),
	//})
	//responseBody := bytes.NewBuffer(postBody)
	//respAccess, err := http.Post(url, "application/json", responseBody)
	//if err != nil {
	//	return "", err
	//}
	//
	//defer respAccess.Body.Close()
	//
	//// Read access body
	//body, err := ioutil.ReadAll(respAccess.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Println(string(body))

	return hex.EncodeToString(buf.Bytes()), err
}
