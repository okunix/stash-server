package crypto

import (
	"crypto/rand"
	"strings"
)

const (
	alphaNumericCharset string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	specialCharset      string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+,-./:;<=>?@[\\]^_{|}~"
)

func RandomBytes(size int) []byte {
	nonce := make([]byte, size)
	rand.Read(nonce)
	return nonce
}

func RandomStringWithCharset(size int, charset string) string {
	b := RandomBytes(size)
	var pass strings.Builder
	pass.Grow(size)
	for _, v := range b {
		pass.WriteByte(charset[int(v)%len(charset)])
	}
	return pass.String()
}

func RandomAlphaNumericString(size int) string {
	return RandomStringWithCharset(size, alphaNumericCharset)
}

func RandomSpecialString(size int) string {
	return RandomStringWithCharset(size, specialCharset)
}
