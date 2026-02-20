package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Cipher interface {
	Encrypt(key, plaintext []byte) ([]byte, error)
	Decrypt(key, ciphertext []byte) ([]byte, error)
}

type AESGCM struct{}

func NewAESGCM() Cipher {
	return AESGCM{}
}

func (st AESGCM) Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

func (st AESGCM) Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := ciphertext[:aesgcm.NonceSize()]
	ciphertext = ciphertext[aesgcm.NonceSize():]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
