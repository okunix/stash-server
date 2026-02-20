package crypto

import (
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

type KDF interface {
	DeriveKey(password, salt []byte) (Key, error)
	Compare(a, b []byte) bool
}

type Key interface {
	String() string
	Bytes() []byte
}

type argon2IDKey struct {
	Argon2ID
	key []byte
}

func (a argon2IDKey) String() string {
	return hex.EncodeToString(a.key)
}
func (a argon2IDKey) Bytes() []byte { return a.key }

type Argon2ID struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type Argon2IDOption func(ar *Argon2ID) error

func WithKeyLen(keyLen uint32) Argon2IDOption {
	return func(ar *Argon2ID) error {
		ar.keyLen = keyLen
		return nil
	}
}

func WithMemory(memory uint32) Argon2IDOption {
	return func(ar *Argon2ID) error {
		ar.memory = memory
		return nil
	}
}

func WithThreads(threads uint8) Argon2IDOption {
	return func(ar *Argon2ID) error {
		ar.threads = threads
		return nil
	}
}

func WithTime(time uint32) Argon2IDOption {
	return func(ar *Argon2ID) error {
		ar.time = time
		return nil
	}
}

func newDefaultArgon2ID() Argon2ID {
	return Argon2ID{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
}

func NewArgon2ID(opts ...Argon2IDOption) (KDF, error) {
	argon2 := newDefaultArgon2ID()
	for _, opt := range opts {
		if err := opt(&argon2); err != nil {
			return nil, err
		}
	}
	return &argon2, nil
}

func (s *Argon2ID) DeriveKey(password, salt []byte) (Key, error) {
	return argon2IDKey{
		Argon2ID: *s,
		key:      argon2.IDKey(password, salt, s.time, s.memory, s.threads, s.keyLen),
	}, nil
}

func (s *Argon2ID) Compare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
