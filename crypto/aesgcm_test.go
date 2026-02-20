package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESGCM(t *testing.T) {
	salt, _ := RandomBytes(32)
	kdf, err := NewArgon2ID(salt)
	assert.Nil(t, err)

	key, err := kdf.DeriveKey([]byte("supeR_stronG_passworD1"))
	assert.Nil(t, err)

	plaintext := []byte("secret message :)")

	cipher := NewAESGCM()
	ciphertext, err := cipher.Encrypt(key.Bytes(), plaintext)
	assert.Nil(t, err)

	t.Log(string(ciphertext))

	decoded, err := cipher.Decrypt(key.Bytes(), ciphertext)
	assert.Nil(t, err)

	t.Log(string(decoded))

	assert.Equal(t, string(plaintext), string(decoded))
}
