package crypto

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESGCM(t *testing.T) {
	kdf, err := NewArgon2ID()
	assert.Nil(t, err)

	key, err := kdf.DeriveKey([]byte("supeR_stronG_passworD1"))
	assert.Nil(t, err)

	plaintext := []byte("secret message :)")

	cipher := AESGCM()
	ciphertext, err := cipher.Encrypt(key.Bytes(), plaintext)
	assert.Nil(t, err)

	t.Log(base64.RawStdEncoding.EncodeToString(ciphertext))

	decoded, err := cipher.Decrypt(key.Bytes(), ciphertext)
	assert.Nil(t, err)

	t.Log(string(decoded))

	assert.Equal(t, string(plaintext), string(decoded))
}
