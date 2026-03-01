package crypto

import "crypto/rand"

func RandomBytes(size int) []byte {
	nonce := make([]byte, size)
	rand.Read(nonce)
	return nonce
}
