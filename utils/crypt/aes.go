package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var JwtEncryptionKey = []byte("0123456789abcdef0123456789abcdef")

func EncryptAES(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		fmt.Println(err)
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(cipherBase64 string, key []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cipherBase64)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	return aesGCM.Open(nil, nonce, cipherText, nil)
}
