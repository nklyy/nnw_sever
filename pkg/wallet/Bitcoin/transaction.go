package Bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"nnw_s/pkg/wallet"
)

var (
	mnemonic = "chair column reveal income inside soul blade concert series syrup ivory bulb"
)

type utxo struct {
	key    *wallet.Key
	hash   string // hash of this utxo. example: f75f1b0810857dff6972767d745a57d0aeb77ae4845ef4f698a7bd50ec34a4b4 (BTC Testnet)
	index  uint32 // index of this wallet at this utxo
	amount int64  // satoshis amount of this utxo
}

func (u *utxo) ToTxIn() (*wire.TxIn, error) {
	hash, err := chainhash.NewHashFromStr(u.hash)
	if err != nil {
		return nil, err
	}

	outPoint := wire.NewOutPoint(hash, u.index)
	return wire.NewTxIn(outPoint, nil, nil), nil
}

func (u *utxo) Signature(tx *wire.MsgTx, index int) error {
	address, err := u.key.AddressBTC()
	if err != nil {
		return err
	}

	addr, err := btcutil.DecodeAddress(address, u.key.Opt.Params)
	if err != nil {
		return err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err
	}

	wif, err := btcutil.NewWIF(u.key.Private, u.key.Opt.Params, false)
	if err != nil {
		return err
	}

	sig, err := txscript.SignatureScript(tx, index, script, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return err
	}

	tx.TxIn[index].SignatureScript = sig
	return nil
}

// receiver is who and how much you want sent coins
type receiver struct {
	net     *chaincfg.Params
	address string // receiver address. example: mxLBntJAV4dF977jJjN6uhT9er9KKeMEgK (BTC Testnet)
	amount  int64  // satoshis amount
}

func (r *receiver) ToTxOut() (*wire.TxOut, error) {
	addr, err := btcutil.DecodeAddress(r.address, r.net)
	if err != nil {
		return nil, err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return wire.NewTxOut(r.amount, script), nil
}

func RunTX() {
	master, err := wallet.NewKey(
		wallet.Mnemonic(mnemonic),
	)
	if err != nil {
		panic(err)
	}

	getWallet, err := master.GetWallet(
		wallet.CoinType(wallet.BtcTestNetType),
		wallet.AddressIndex(1),
	)
	if err != nil {
		panic(err)
	}

	key := getWallet.GetKey()

	senders := []*utxo{
		{key: key, hash: "4ef4e638013922605091d53e2164864ef602ce665926697f3082168ccb2eb16d", index: 0, amount: 990000},
	}

	receivers := []*receiver{
		{net: key.Opt.Params, address: "mz5otQWaaPmWEydGw9gjfaHxSKhHX6Fgj5", amount: 990000},
	}

	tx := wire.NewMsgTx(wire.TxVersion)

	for _, s := range senders {
		in, err := s.ToTxIn()
		if err != nil {
			panic(err)
		}

		tx.AddTxIn(in)
	}

	for _, r := range receivers {
		out, err := r.ToTxOut()
		if err != nil {
			panic(err)
		}

		tx.AddTxOut(out)
	}

	for i, s := range senders {
		if err := s.Signature(tx, i); err != nil {
			panic(err)
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	fmt.Println("Transaction:", hex.EncodeToString(buf.Bytes()))
}
