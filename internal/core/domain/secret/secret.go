package secret

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/crypto"
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

func (s *Secret) Seal(cipher crypto.Cipher) ([]byte, error) {
	data, _ := json.Marshal(s.Data)
	return cipher.Encrypt(s.MasterKey, data)
}

func NewFromCipher(cipher crypto.Cipher, key, ciphertext []byte) (Secret, error) {
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
