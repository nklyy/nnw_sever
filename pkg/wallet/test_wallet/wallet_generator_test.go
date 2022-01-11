package test_wallet

import (
	"fmt"
	"testing"
)

func TestGenerateWallet(t *testing.T) {
	secretPassphrase := "secret"
	generated := Generate(secretPassphrase)

	// Display wallet and keys
	fmt.Println("Mnemonic: ", generated.Mnemonic)
	fmt.Println("Master private key: ", generated.MasterKey)
	fmt.Println("Master public key: ", generated.PublicKey)

	mnemonicEncrypted := Encrypt([]byte(generated.Mnemonic), secretPassphrase)
	mnemonicDecrypted := Decrypt(mnemonicEncrypted, secretPassphrase)

	fmt.Println("Decrypted Mnemonic:", string(mnemonicDecrypted))

	if string(mnemonicDecrypted) != generated.Mnemonic {
		t.Error("mnemonic does not valid!")
	}

}
