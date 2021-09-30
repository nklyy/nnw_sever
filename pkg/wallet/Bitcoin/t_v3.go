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
)

func RunTransactionV3() error {
	privWif := "cUNyvH5eyeVsw7eXsGRqis3joYBveXBQQRMF8pybWdjSNwszUTDP"
	txHash := "2ad8223b1869853d37ce0df681d5840f306de8b56b8bd401bff137d412946e63"
	destination := "mrMb6DZLJG9fXafGAB53Gs9GDS4CNerP5t"
	amount := int64(10000)
	txFee := int64(200)
	sourceUTXOIndex := uint32(1)
	chainParams := &chaincfg.TestNet3Params

	decodedWif, err := btcutil.DecodeWIF(privWif)
	if err != nil {
		return err
	}

	fmt.Printf("Decoded WIF: %v\n", decodedWif) // Decoded WIF: cTujmQgVdGYzmZEfhq5gVDpd2EAHF1sZahPmkDnRHmPDEVRYz6eo

	addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)
	if err != nil {
		return err
	}

	sourceUTXOHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return err
	}

	fmt.Printf("UTXO hash: %s\n", sourceUTXOHash) // utxo hash: c6950f355835c361dce2e9d6eb511cf56972b67cb34dad5d1fd9f9bc796711a5

	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(destination, chainParams)
	if err != nil {
		return err
	}

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	if err != nil {
		return err
	}

	fmt.Printf("Source Address: %s\n", sourceAddress) // Source Address: mqJ8FALtYnxvLgwTUWQ2shNkdiLuU7tkPR

	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return err
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return err
	}

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	redeemTx.AddTxIn(sourceTxIn)
	redeemTxOut := wire.NewTxOut((amount - txFee), destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTxOut.PkScript, txscript.SigHashAll, decodedWif.PrivKey, true)
	if err != nil {
		return err
	}

	redeemTx.TxIn[0].SignatureScript = sigScript
	fmt.Printf("Signature Script: %v\n", hex.EncodeToString(sigScript)) // Signature Script: 473...b67

	// validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTxOut.PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return err
	}

	if err := vm.Execute(); err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0, redeemTx.SerializeSize()))
	redeemTx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes())) // redeem Tx: 01000000011efc...5bb88ac00000000

	//bcy := gobcy.API{"55f0c359f95b4bc5a1c6e949c8c74731", "btc", "test3"}
	//skel, err := bcy.PushTX(hex.EncodeToString(buf.Bytes()))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%+v\n", skel)

	return nil
}
