package service

import (
	"crypto/aes"
	"crypto/cipher"
)

// VerifyService check the code that comes back is a valid login
func VerifyService(body string) (string, error) {
	return "login failed", nil
}

func decryptData(code string, id int) string {
	key := []byte(code)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	data := RetrieveData(id)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}
