package crypto

import "crypto/rand"

func RandomBytes(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	return nonce, err
}
