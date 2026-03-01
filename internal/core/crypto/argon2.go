package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
)

var (
	argon2EncodingRegex = regexp.MustCompile(
		`^\$argon2id\$v=(?P<v>\d+)\$m=(?P<d>\d+),t=(?P<t>\d+),p=(?P<p>\d+)\$(?P<salt>\S+?)\$(?P<hash>\S+)$`,
	)
)

type KDF interface {
	DeriveKey(password []byte) (Key, error)
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
	saltb64 := base64.RawStdEncoding.EncodeToString(a.salt)
	keyb64 := base64.RawStdEncoding.EncodeToString(a.key)
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.memory,
		a.time,
		a.threads,
		saltb64,
		keyb64,
	)
	return encoded
}

func (a argon2IDKey) Bytes() []byte {
	return a.key
}

type Argon2ID struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	salt    []byte
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

func WithSalt(salt []byte) Argon2IDOption {
	return func(ar *Argon2ID) error {
		ar.salt = salt
		return nil
	}
}

func newDefaultArgon2ID() Argon2ID {
	randomSalt := make([]byte, 32)
	rand.Read(randomSalt)
	return Argon2ID{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
		salt:    randomSalt,
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

func NewArgon2IDFromString(s string) (kdf KDF, hash []byte, err error) {
	matches := argon2EncodingRegex.FindStringSubmatch(s)
	salt, _ := base64.RawStdEncoding.DecodeString(matches[5])
	hash, _ = base64.RawStdEncoding.DecodeString(matches[6])
	memory, _ := strconv.ParseUint(matches[2], 10, 32)
	time, _ := strconv.ParseUint(matches[3], 10, 32)
	threads, _ := strconv.ParseUint(matches[4], 10, 8)
	kdf, err = NewArgon2ID(
		WithSalt(salt),
		WithMemory(uint32(memory)),
		WithThreads(uint8(threads)),
		WithTime(uint32(time)),
	)
	return
}

func (s *Argon2ID) DeriveKey(password []byte) (Key, error) {
	return argon2IDKey{
		Argon2ID: *s,
		key:      argon2.IDKey(password, s.salt, s.time, s.memory, s.threads, s.keyLen),
	}, nil
}

func (s *Argon2ID) Compare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
