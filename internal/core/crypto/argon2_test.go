package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2(t *testing.T) {
	password := []byte("helloworld")
	salt := []byte("aaaaaaaa")
	kdf, err := NewArgon2ID(
		WithSalt(salt),
		WithMemory(64*1024),
		WithKeyLen(32),
		WithTime(1),
		WithThreads(4),
	)
	assert.Nil(t, err)
	key1, err := kdf.DeriveKey(password)
	assert.Nil(t, err)

	key2, err := kdf.DeriveKey(password)
	assert.Nil(t, err)
	assert.Equal(t, key1.String(), key2.String())

	t.Log(key1)
	t.Log(key2)

	kdf, hash, err := NewArgon2IDFromString(key1.String())
	key3, err := kdf.DeriveKey(password)
	assert.Nil(t, err)

	assert.Equal(t, key2, key3)
	assert.Equal(t, string(hash), string(key2.Bytes()))

	header := string(key1.Salt())
	t.Log(header)
	kdf, err = NewArgon2ID(WithHeader(header))
	assert.Nil(t, err)
	key4, err := kdf.DeriveKey(password)
	assert.Nil(t, err)
	assert.Equal(t, key1.String(), key4.String())
	t.Log(key1.String())
	t.Log(key4.String())
}
