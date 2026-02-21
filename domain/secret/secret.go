package secret

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/crypto"
)

// in-memory data structure
type Secret struct {
	mu         sync.RWMutex
	MasterKey  []byte            `json:"master_key"`
	Data       map[string]string `json:"data"`
	UnlockedAt time.Time         `json:"unlocked_at"`
}

type AddSecretParams struct {
	StashID      uuid.UUID
	MaintainerID uuid.UUID
	MasterKey    []byte
	Data         map[string]string
}

type GetSecretParams struct {
	StashID uuid.UUID
}

type RemoveSecretParams struct {
	StashID uuid.UUID
}

type Repository interface {
	AddSecret(ctx context.Context, params AddSecretParams) error
	RemoveSecret(ctx context.Context, params RemoveSecretParams) error
	GetSecret(ctx context.Context, params GetSecretParams) (*Secret, error)
}

func (s *Secret) Seal() ([]byte, error) {
	cipher := crypto.AESGCM()
	data, _ := json.Marshal(s.Data)
	return cipher.Encrypt(s.MasterKey, data)
}

func NewFromCipher(key, ciphertext []byte) (Secret, error) {
	cipher := crypto.AESGCM()
	plaintext, err := cipher.Decrypt(key, ciphertext)
	if err != nil {
		return Secret{}, err
	}
	var data map[string]string
	err = json.Unmarshal([]byte(plaintext), &data)
	return Secret{
		MasterKey:  key,
		Data:       data,
		UnlockedAt: time.Now(),
	}, err
}

func (s *Secret) AddEntry(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
}

func (s *Secret) RemoveEntry(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data, key)
}

func (s *Secret) GetEntry(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Data[key]
}
