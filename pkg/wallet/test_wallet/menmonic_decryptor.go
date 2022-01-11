package test_wallet

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(data []byte, passphrase string) []byte {
	// Create passphrase hash key
	key := []byte(createHash(passphrase))

	// Create a new block cipher based on the hashed passphrase
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Wrap block in (GCM) with a standard nonce length.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Separating the nonce and the encrypted data.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext
}
