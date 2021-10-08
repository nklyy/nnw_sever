package feature

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/wire"
	"log"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

//// UTXO ...
//type UTXO struct {
//	Hash      string
//	TxIndex   int
//	Amount    *big.Int
//	Spendable bool
//	PKScript  []byte
//}

func marshalUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
	// same strategy as bitcoin core
	// from: https://medium.com/@lopp/the-challenges-of-optimizing-unspent-output-selection-a3e5d05d13ef
	// 1. sort the UTXOs from smallest to largest amounts
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount.Cmp(utxos[j].Amount) == -1
	})

	// 2. search for exact match
	for idx := range utxos {
		exactTxSize := calculateTotalTxBytes(1, 2)
		totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
		totalTxAmount := new(big.Int).Add(totalFee, amount)

		switch utxos[idx].Amount.Cmp(totalTxAmount) {
		case 0:
			var resp []*UTXO
			resp = append(resp, utxos[idx])
			// TODO: store these in the DB to be sure they aren't being claimed??
			return resp, sumUTXOs(resp), nil

		case 1:
			break
		}
	}

	// 3. calculate the sum of all UTXOs smaller than amount
	sumSmall := big.NewInt(0)
	var sumSmallUTXOs []*UTXO
	for idx := range utxos {
		switch utxos[idx].Amount.Cmp(amount) {
		case -1:
			_ = sumSmall.Add(sumSmall, utxos[idx].Amount)
			sumSmallUTXOs = append(sumSmallUTXOs, utxos[idx])

		default:
			break
		}
	}

	exactTxSize := calculateTotalTxBytes(len(sumSmallUTXOs), 2)
	totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
	totalTxAmount := new(big.Int).Add(totalFee, amount)

	switch sumSmall.Cmp(totalTxAmount) {
	case 0:
		return sumSmallUTXOs, sumUTXOs(sumSmallUTXOs), nil

	case -1:
		for idx := range utxos {
			exactTxSize := calculateTotalTxBytes(1, 2)
			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
			totalTxAmount := new(big.Int).Add(totalFee, amount)
			if utxos[idx].Amount.Cmp(totalTxAmount) == 1 {
				var resp []*UTXO
				resp = append(resp, utxos[idx])
				return resp, sumUTXOs(resp), nil
			}
		}

		// should reach here if not enought UXOs
		log.Fatal("not enough UTXOs to meet target amount")

	case 1:
		return roundRobinSelectUTXOs(sumSmallUTXOs, amount, feeRate)

	default:
		log.Fatal("unknown comparison")
	}

	return nil, nil, nil
}

func roundRobinSelectUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
	var possibilities [][]*UTXO
	lenInput := len(utxos)
	log.Printf("round robin select; lenInput: %v", lenInput)
	if lenInput == 0 {
		log.Fatal("expected utxos size to be greater than 0")
	}

	for i := 0; i < 1000; i++ {
		selectedIdxs := make(map[int]bool)
		var sum *big.Int
		var possibility []*UTXO
		for {
			for {
				rand.Seed(time.Now().Unix())
				tmp := 0
				if lenInput > 1 {
					tmp = rand.Intn(lenInput - 1)
				}

				if !selectedIdxs[tmp] {
					selectedIdxs[tmp] = true
					_ = sum.Add(sum, utxos[tmp].Amount)
					possibility = append(possibility, utxos[tmp])

					break
				}
			}

			exactTxSize := calculateTotalTxBytes(len(possibility), 2)
			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
			totalTxAmount := new(big.Int).Add(totalFee, amount)

			if sum.Cmp(totalTxAmount) == 0 {
				return possibility, sum, nil
			}

			if sum.Cmp(totalTxAmount) == 1 {
				possibilities = append(possibilities, possibility)
				break
			}
		}
	}

	if len(possibilities) < 1 {
		return nil, nil, errors.New("no possible utxo combos")
	}

	smallestLen := len(possibilities[0])
	smallestIdx := 0

	for idx := 1; idx < len(possibilities); idx++ {
		l := len(possibilities[idx])
		if l < smallestLen {
			smallestLen = l
			smallestIdx = idx
		}
	}

	return possibilities[smallestIdx], sumUTXOs(possibilities[smallestIdx]), nil
}

func sumUTXOs(utxos []*UTXO) *big.Int {
	sum := big.NewInt(0)
	for idx := range utxos {
		sum = sum.Add(sum, utxos[idx].Amount)
	}

	return sum
}

// https://bitcoin.stackexchange.com/questions/1195/how-to-calculate-transaction-size-before-sending-legacy-non-segwit-p2pkh-p2sh
func calculateTotalTxBytes(txInLength, txOutLength int) int {
	return txInLength*180 + txOutLength*34 + 10 + txInLength
}

func decodeRawTx(rawTx string) (*wire.MsgTx, error) {
	raw, err := hex.DecodeString(rawTx)
	if err != nil {
		log.Printf("err decoding raw tx; err: %v", err)
		return nil, err
	}

	var version int32 = 2
	if rawTx[:8] == "01000000" {
		version = 1
	}
	log.Printf("version: %d", version)

	r := bytes.NewReader(raw)
	tmpTx := wire.NewMsgTx(version)

	err = tmpTx.BtcDecode(r, uint32(version), wire.BaseEncoding)
	if err != nil {
		log.Printf("could not decode raw tx; err: %v", err)
		return nil, err
	}

	return tmpTx, nil
}
