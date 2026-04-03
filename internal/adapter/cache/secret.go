package cache

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/core/domain/secret"
	"github.com/okunix/stash-server/internal/core/ports"
)

type secretRepository struct {
	cache       *Cache
	userMapping *Cache
}

func NewSecretRepository(cache, userMapping *Cache) ports.SecretRepository {
	return &secretRepository{cache: cache, userMapping: userMapping}
}

func (s *secretRepository) GetSecretByStashID(
	ctx context.Context,
	stashID uuid.UUID,
) (*secret.Secret, error) {
	st, ok := s.cache.Get(stashID.String())
	if !ok {
		return nil, errors.New("stash is locked")
	}
	secretData := st.(*secret.Secret)
	return secretData, nil
}

func (s *secretRepository) AddSecret(
	ctx context.Context,
	params secret.AddSecretParams,
) (*secret.Secret, error) {
	newSecret := secret.Secret{
		MasterKey:    params.MasterKey,
		Data:         params.Data,
		MaintainerID: params.MaintainerID,
		UnlockedAt:   time.Now(),
	}
	s.cache.Set(params.StashID.String(), &newSecret)

	newMap := []uuid.UUID{params.StashID}
	userMap, ok := s.userMapping.Get(params.MaintainerID.String())
	if ok {
		newMap = userMap.([]uuid.UUID)
		newMap = append(newMap, params.StashID)
	}
	s.userMapping.Set(params.MaintainerID.String(), newMap)

	return &newSecret, nil
}

func (s *secretRepository) UpdateSecret(
	ctx context.Context,
	stashID uuid.UUID,
	sec *secret.Secret,
) error {
	s.cache.Set(stashID.String(), sec)
	return nil
}

func (s *secretRepository) RemoveSecretByStashID(
	ctx context.Context,
	stashID uuid.UUID,
) (*secret.Secret, error) {
	sec, err := s.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return nil, err
	}
	s.cache.Delete(stashID.String())
	return sec, nil
}

func (s *secretRepository) ListSecrets(
	ctx context.Context,
	maintainerID uuid.UUID,
) ([]*secret.Secret, error) {
	userMap, ok := s.userMapping.Get(maintainerID.String())
	if !ok {
		return []*secret.Secret{}, nil
	}
	stashIDs := userMap.([]uuid.UUID)
	secrets := []*secret.Secret{}
	for _, v := range stashIDs {
		sec, err := s.GetSecretByStashID(ctx, v)
		if err != nil {
			continue
		}
		secrets = append(secrets, sec)
	}
	return secrets, nil
}
