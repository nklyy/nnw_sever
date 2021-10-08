package feature

//
//import (
//	"bytes"
//	"encoding/hex"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/btcsuite/btcd/chaincfg"
//	"github.com/btcsuite/btcd/wire"
//	"io/ioutil"
//	"log"
//	"math/big"
//	"math/rand"
//	"net/http"
//	"sort"
//	"time"
//)
//
//// UTXO ...
//type UTXO struct {
//	Hash      string
//	TxIndex   int
//	Amount    *big.Int
//	Spendable bool
//	PKScript  []byte
//}
//
////func sendMsg(req, res interface{}) {
////	//serverAddr := "electrum.qtornado.com:50002" // mainnet
////	//serverAddr := "testnet.qtornado.com:51002" // testnet
////	//serverAddr := "testnet1.bauerj.eu:50002" // testnet
////	//serverAddr := "testnet.hsmiths.com:53012" // testnet
////	//serverAddr := "testnet.aranguren.org:51001" // testnet
////	serverAddr := "http://127.0.0.1:8332" // testnet
////
////	fmt.Printf("dialing to server: %s\n", serverAddr)
////	//conn, err := tls.Dial("tcp", serverAddr, &tls.Config{
////	//	InsecureSkipVerify: false,
////	//})
////	conn, err := net.Dial("tcp", serverAddr)
////	if err != nil {
////		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa")
////		log.Fatal(err)
////	}
////
////	defer conn.Close()
////	fmt.Printf("client connected to: %s\n", conn.RemoteAddr())
////
////	reqMsgBytes, err := json.Marshal(req)
////	if err != nil {
////		log.Fatal(err)
////	}
////
////	reqMsg := fmt.Sprintf("%s\n", string(reqMsgBytes))
////	fmt.Printf("writing message: %s", reqMsg)
////	_, err = io.WriteString(conn, reqMsg)
////	if err != nil {
////		log.Fatal(err)
////	}
////
////	var (
////		i        int
////		readSize int = 1024
////		respData []byte
////	)
////
////	for {
////		fmt.Println("reading response...")
////		respBytes := make([]byte, readSize)
////		n, err := conn.Read(respBytes)
////		if err != nil {
////			if err != io.EOF {
////				log.Fatal(err)
////			}
////		}
////
////		fmt.Printf("reading: %q (%d bytes)\n", string(respBytes[:n]), n)
////
////		respData = append(respData, respBytes[:n]...)
////		i += n
////
////		if n < readSize {
////			break
////		}
////	}
////
////	err = json.Unmarshal(respData[:i], &res)
////	if err != nil {
////		log.Fatal(err)
////	}
////}
//
//func StartProcess() {
//	//chainParams := &chaincfg.MainNetParams
//	chainParams := &chaincfg.TestNet3Params
//
//	amountToSend := big.NewInt(50000) // amount to send in satoshis (0.01 btc)
//
//	feeRate, err := GetCurrentFeeRate()
//	log.Printf("current fee rate: %v", feeRate)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fromWalletPublicAddress := "mq6Qd7JJKsgBYkMFsGCk24MHMxUkuyTnkU"
//
//	log.Printf("from wallet public address: %s", fromWalletPublicAddress)
//
//	unspentTXOs, err := ListUnspentTXOs(fromWalletPublicAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//unspentTXOs, UTXOsAmount, err := marshalUTXOs(unspentTXOs, amountToSend, feeRate)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//// prepare unspent transaction outputs with its privatekey.
//	//log.Println("unspent UTXOs", unspentTXOs, UTXOsAmount)
//	//
//	//tx := wire.NewMsgTx(wire.TxVersion)
//	//
//	//var sourceUTXOs []*UTXO
//	//// prepare tx ins
//	//for idx := range unspentTXOs {
//	//	hashStr := unspentTXOs[idx].Hash
//	//
//	//	sourceUTXOHash, err := chainhash.NewHashFromStr(hashStr)
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//
//	//	sourceUTXOIndex := uint32(unspentTXOs[idx].TxIndex)
//	//	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
//	//	sourceUTXOs = append(sourceUTXOs, unspentTXOs[idx])
//	//	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
//	//
//	//	tx.AddTxIn(sourceTxIn)
//	//}
//	//
//	//// calculate fees
//	//txByteSize := big.NewInt(int64(len(tx.TxIn)*180 + len(tx.TxOut)*34 + 10 + len(tx.TxIn)))
//	//totalFee := new(big.Int).Mul(feeRate, txByteSize)
//	//log.Printf("total fee: %s", totalFee)
//	//
//	//// calculate the change
//	//change := new(big.Int).Set(UTXOsAmount)
//	//change = new(big.Int).Sub(change, amountToSend)
//	//change = new(big.Int).Sub(change, totalFee)
//	//if change.Cmp(big.NewInt(0)) == -1 {
//	//	log.Fatal(err)
//	//}
//	//
//	//destinationAddress := "mmfbzo2533SFa34ErmYNY4RdVtfw5XYK1u"
//	//
//	//// create the tx outs
//	//destAddress, err := btcutil.DecodeAddress(destinationAddress, chainParams)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//destScript, err := txscript.PayToAddrScript(destAddress)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//// tx out to send btc to user
//	//destOutput := wire.NewTxOut(amountToSend.Int64(), destScript)
//	//tx.AddTxOut(destOutput)
//	//
//	//// our change address
//	//changeSendToAddress, err := btcutil.DecodeAddress(fromWalletPublicAddress, chainParams)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//changeSendToScript, err := txscript.PayToAddrScript(changeSendToAddress)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//// tx out to send change back to us
//	//changeOutput := wire.NewTxOut(change.Int64(), changeSendToScript)
//	//tx.AddTxOut(changeOutput)
//	//
//	//privWif := "cPRZfnSdhrLvetS9KySaxdqD99yoy1mD3tHhDaMRDqM1gdWf36KD"
//	//
//	//decodedWif, err := btcutil.DecodeWIF(privWif)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//fmt.Printf("Source Address: %s\n", sourceAddress) // Source Address: mgjHgKi1g6qLFBM1gQwuMjjVBGMJdrs9pP
//	//
//	//sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//for i := range sourceUTXOs {
//	//	sigScript, err := txscript.SignatureScript(tx, i, sourcePkScript, txscript.SigHashAll, decodedWif.PrivKey, true)
//	//	if err != nil {
//	//		log.Fatalf("could not generate pubSig; err: %v", err)
//	//	}
//	//	tx.TxIn[i].SignatureScript = sigScript
//	//}
//	//
//	//buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
//	//tx.Serialize(buf)
//	//
//	//fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes()))
//	//
//	//t := hex.EncodeToString(buf.Bytes())
//	//txHash, err := SendTX(t)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//fmt.Printf("tx hash: %s\n", txHash) // 1d8f70dfc8b90bff672ee663a7cc811c4e88e98c6895dc93aa9f73202bb7809b
//}
//
//func marshalUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
//	// same strategy as bitcoin core
//	// from: https://medium.com/@lopp/the-challenges-of-optimizing-unspent-output-selection-a3e5d05d13ef
//	// 1. sort the UTXOs from smallest to largest amounts
//	sort.Slice(utxos, func(i, j int) bool {
//		return utxos[i].Amount.Cmp(utxos[j].Amount) == -1
//	})
//
//	// 2. search for exact match
//	for idx := range utxos {
//		exactTxSize := calculateTotalTxBytes(1, 2)
//		totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
//		totalTxAmount := new(big.Int).Add(totalFee, amount)
//
//		switch utxos[idx].Amount.Cmp(totalTxAmount) {
//		case 0:
//			var resp []*UTXO
//			resp = append(resp, utxos[idx])
//			// TODO: store these in the DB to be sure they aren't being claimed??
//			return resp, sumUTXOs(resp), nil
//
//		case 1:
//			break
//		}
//	}
//
//	// 3. calculate the sum of all UTXOs smaller than amount
//	sumSmall := big.NewInt(0)
//	var sumSmallUTXOs []*UTXO
//	for idx := range utxos {
//		switch utxos[idx].Amount.Cmp(amount) {
//		case -1:
//			_ = sumSmall.Add(sumSmall, utxos[idx].Amount)
//			sumSmallUTXOs = append(sumSmallUTXOs, utxos[idx])
//
//		default:
//			break
//		}
//	}
//
//	exactTxSize := calculateTotalTxBytes(len(sumSmallUTXOs), 2)
//	totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
//	totalTxAmount := new(big.Int).Add(totalFee, amount)
//
//	switch sumSmall.Cmp(totalTxAmount) {
//	case 0:
//		return sumSmallUTXOs, sumUTXOs(sumSmallUTXOs), nil
//
//	case -1:
//		for idx := range utxos {
//			exactTxSize := calculateTotalTxBytes(1, 2)
//			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
//			totalTxAmount := new(big.Int).Add(totalFee, amount)
//			if utxos[idx].Amount.Cmp(totalTxAmount) == 1 {
//				var resp []*UTXO
//				resp = append(resp, utxos[idx])
//				return resp, sumUTXOs(resp), nil
//			}
//		}
//
//		// should reach here if not enought UXOs
//		log.Fatal("not enough UTXOs to meet target amount")
//
//	case 1:
//		return roundRobinSelectUTXOs(sumSmallUTXOs, amount, feeRate)
//
//	default:
//		log.Fatal("unknown comparison")
//	}
//
//	return nil, nil, nil
//}
//
//func roundRobinSelectUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
//	var possibilities [][]*UTXO
//	lenInput := len(utxos)
//	log.Printf("round robin select; lenInput: %v", lenInput)
//	if lenInput == 0 {
//		log.Fatal("expected utxos size to be greater than 0")
//	}
//
//	for i := 0; i < 1000; i++ {
//		selectedIdxs := make(map[int]bool)
//		var sum *big.Int
//		var possibility []*UTXO
//		for {
//			for {
//				rand.Seed(time.Now().Unix())
//				tmp := 0
//				if lenInput > 1 {
//					tmp = rand.Intn(lenInput - 1)
//				}
//
//				if !selectedIdxs[tmp] {
//					selectedIdxs[tmp] = true
//					_ = sum.Add(sum, utxos[tmp].Amount)
//					possibility = append(possibility, utxos[tmp])
//
//					break
//				}
//			}
//
//			exactTxSize := calculateTotalTxBytes(len(possibility), 2)
//			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
//			totalTxAmount := new(big.Int).Add(totalFee, amount)
//
//			if sum.Cmp(totalTxAmount) == 0 {
//				return possibility, sum, nil
//			}
//
//			if sum.Cmp(totalTxAmount) == 1 {
//				possibilities = append(possibilities, possibility)
//				break
//			}
//		}
//	}
//
//	if len(possibilities) < 1 {
//		return nil, nil, errors.New("no possible utxo combos")
//	}
//
//	smallestLen := len(possibilities[0])
//	smallestIdx := 0
//
//	for idx := 1; idx < len(possibilities); idx++ {
//		l := len(possibilities[idx])
//		if l < smallestLen {
//			smallestLen = l
//			smallestIdx = idx
//		}
//	}
//
//	return possibilities[smallestIdx], sumUTXOs(possibilities[smallestIdx]), nil
//}
//
//func sumUTXOs(utxos []*UTXO) *big.Int {
//	sum := big.NewInt(0)
//	for idx := range utxos {
//		sum = sum.Add(sum, utxos[idx].Amount)
//	}
//
//	return sum
//}
//
//// https://bitcoin.stackexchange.com/questions/1195/how-to-calculate-transaction-size-before-sending-legacy-non-segwit-p2pkh-p2sh
//func calculateTotalTxBytes(txInLength, txOutLength int) int {
//	return txInLength*180 + txOutLength*34 + 10 + txInLength
//}
//
//func decodeRawTx(rawTx string) (*wire.MsgTx, error) {
//	raw, err := hex.DecodeString(rawTx)
//	if err != nil {
//		log.Printf("err decoding raw tx; err: %v", err)
//		return nil, err
//	}
//
//	var version int32 = 2
//	if rawTx[:8] == "01000000" {
//		version = 1
//	}
//	log.Printf("version: %d", version)
//
//	r := bytes.NewReader(raw)
//	tmpTx := wire.NewMsgTx(version)
//
//	err = tmpTx.BtcDecode(r, uint32(version), wire.BaseEncoding)
//	if err != nil {
//		log.Printf("could not decode raw tx; err: %v", err)
//		return nil, err
//	}
//
//	return tmpTx, nil
//}
//
//// ListUnspentTXOs lists all UTXOs for an address
//func ListUnspentTXOs(address string) ([]*UTXO, error) {
//	req := struct {
//		ID     int      `json:"id"`
//		Method string   `json:"method"`
//		Params []string `json:"params"`
//	}{
//		ID:     1,
//		Method: "blockchain.address.listunspent",
//		Params: []string{address},
//	}
//
//	msg := struct {
//		JSONRPC string `json:"jsonrpc,omitempty"`
//		ID      int    `json:"id"`
//		Result  []struct {
//			TXHash     string   `json:"tx_hash"`
//			TXPosition uint64   `json:"tx_pos"`
//			Value      *big.Int `json:"value"`
//			Height     uint64   `json:"height"`
//		} `json:"result"`
//	}{}
//
//	var MaxTries = 5
//	for try := 0; try < MaxTries; try++ {
//		sendMsgNew(req, &msg)
//
//		var utxos []*UTXO
//		for idx := range msg.Result {
//			utxos = append(utxos, &UTXO{
//				Hash:      msg.Result[idx].TXHash,
//				TxIndex:   int(msg.Result[idx].TXPosition),
//				Amount:    msg.Result[idx].Value,
//				Spendable: true,
//			})
//		}
//
//		return utxos, nil
//	}
//
//	log.Printf("could not get utxos")
//	return nil, errors.New("could not get utxos")
//}
//
//// GetRawTransaction gets raw transaction data given transaction ID (hash)
////func GetRawTransaction(txHash string) ([]byte, error) {
////	req := struct {
////		ID     int      `json:"id"`
////		Method string   `json:"method"`
////		Params []string `json:"params"`
////	}{
////		ID:     1,
////		Method: "blockchain.transaction.get",
////		Params: []string{txHash},
////	}
////
////	msg := struct {
////		JSONRPC string `json:"jsonrpc,omitempty"`
////		ID      int    `json:"id"`
////		Result  string `json:"result"`
////	}{}
////
////	var MaxTries = 5
////	for try := 0; try < MaxTries; try++ {
////		sendMsg(req, &msg)
////
////		b, err := hex.DecodeString(msg.Result)
////		if err != nil {
////			log.Printf("could not decode tx raw data to bytes; err: %v", err)
////			return nil, err
////		}
////
////		return b, nil
////	}
////
////	log.Print("could not get transaction info")
////	return nil, errors.New("could not get transaction info")
////}
////
////// GetTransaction gets transaction data given transaction ID (hash)
////func GetTransaction(txHash string) (*wire.MsgTx, error) {
////	rawTx, err := GetRawTransaction(txHash)
////	if err != nil {
////		log.Printf("err getting raw tx; err: %v", err)
////		return nil, err
////	}
////
////	fmt.Println("RAW", hex.EncodeToString(rawTx))
////
////	tx, err := decodeRawTx(hex.EncodeToString(rawTx))
////	if err != nil {
////		log.Printf("err parsing raw tx; err: %v", err)
////		return nil, err
////	}
////
////	return tx, nil
////}
////
////// SendTX sends a transaction on the wire
////func SendTX(tx string) (string, error) {
////	req := struct {
////		ID     int      `json:"id"`
////		Method string   `json:"method"`
////		Params []string `json:"params"`
////	}{
////		ID:     1,
////		Method: "blockchain.transaction.broadcast",
////		Params: []string{tx},
////	}
////
////	msg := struct {
////		JSONRPC string `json:"jsonrpc,omitempty"`
////		ID      int    `json:"id"`
////		Result  string `json:"result"`
////	}{}
////
////	log.Print("attempting to send bitcoin tx")
////	var MaxTries = 5
////	for try := 0; try < MaxTries; try++ {
////		sendMsg(req, &msg)
////
////		return msg.Result, nil
////	}
////
////	log.Print("could not broadcast tx")
////	return "", errors.New("could not broadcast tx")
////}
