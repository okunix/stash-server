package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Cipher interface {
	Encrypt(key, nonce, data []byte) ([]byte, error)
	Decrypt(key, nonce, data []byte) ([]byte, error)
}

type AESGCM struct{}

func Encrypt(data string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, []byte(data), nil)
	return ciphertext, nil
}

func RandomBytes(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	return nonce, err
}
