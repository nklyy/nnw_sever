package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Hash key, then return the hash as a hexadecimal value.
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	// Create a new block cipher based on the hashed passphrase
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))

	// Wrap block in (GCM) with a standard nonce length.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Create a nonce.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Append the nonce to the encrypted data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}
