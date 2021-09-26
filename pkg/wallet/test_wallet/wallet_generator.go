package test_wallet

import (
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Generator struct {
	Mnemonic  string
	MasterKey *bip32.Key
	PublicKey *bip32.Key
}

func Generate(secretPassphrase string) Generator {
	// Generate a wallet for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Generate a Bip32 HD wallet for the wallet and a user supplied password
	seed := bip39.NewSeed(mnemonic, secretPassphrase)

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	return Generator{
		Mnemonic:  mnemonic,
		MasterKey: masterKey,
		PublicKey: publicKey,
	}
}
