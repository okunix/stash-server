package cache

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/secret"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type secretRepository struct {
	cache *Cache
}

func NewSecretRepository(cache *Cache) ports.SecretRepository {
	return &secretRepository{cache: cache}
}

func (s *secretRepository) GetSecretByStashID(
	ctx context.Context,
	stashID uuid.UUID,
) (*secret.Secret, error) {
	st, ok := s.cache.Get(stashID.String())
	if !ok {
		return nil, errors.New("secret not found")
	}
	secretData := st.(*secret.Secret)
	return secretData, nil
}

func (s *secretRepository) AddSecret(ctx context.Context, params secret.AddSecretParams) error {
	newSecret := secret.Secret{
		MasterKey:  params.MasterKey,
		Data:       params.Data,
		UnlockedAt: time.Now(),
	}
	s.cache.Set(params.StashID.String(), &newSecret)
	return nil
}

func (s *secretRepository) RemoveSecretByStashID(ctx context.Context, stashID uuid.UUID) error {
	s.cache.Delete(stashID.String())
	return nil
}
