package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func CreateWallet(coinType uint32, mnemonic string) (*hdkeychain.ExtendedKey, error) {
	seed := bip39.NewSeed(mnemonic, "")

	// Generate a new master node using the seed.
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H
	purposePath, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H
	coinTypePath, err := purposePath.Derive(hdkeychain.HardenedKeyStart + coinType)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H
	accountPath, err := coinTypePath.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H/0
	changePath, err := accountPath.Derive(0)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H/0/0
	indexPath, err := changePath.Derive(0)
	if err != nil {
		return nil, err
	}

	return indexPath, nil
}

func ToBTCWallet(key *hdkeychain.ExtendedKey) (*BTCWallet, error) {
	address, err := key.Address(&chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	ECPKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	privKey, _ := btcec.PrivKeyFromBytes(ECPKey.ToECDSA().PublicKey.Curve, crypto.FromECDSA(ECPKey.ToECDSA()))
	btcwif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)

	return &BTCWallet{
		Address:    address.String(),
		PrivateKey: btcwif.String(),
	}, nil
}

func ToETHWallet(key *hdkeychain.ExtendedKey) (*ETHWallet, error) {
	ECPKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	keyBytes := crypto.FromECDSA(ECPKey.ToECDSA())

	address := crypto.PubkeyToAddress(ECPKey.ToECDSA().PublicKey).Hex()
	privateKey := hexutil.Encode(keyBytes)

	return &ETHWallet{
		Address:    address,
		PrivateKey: privateKey,
	}, nil
}
