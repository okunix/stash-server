package repository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/stash-password-manager/stash-server/cache"
	"gitlab.com/stash-password-manager/stash-server/domain/secret"
)

type SecretRepositoryImpl struct {
	cache *cache.Cache
}

func NewSecretRepository(cache *cache.Cache) secret.Repository {
	return &SecretRepositoryImpl{cache: cache}
}

func (s *SecretRepositoryImpl) GetSecret(
	ctx context.Context,
	params secret.GetSecretParams,
) (*secret.Secret, error) {
	st, ok := s.cache.Get(params.StashID.String())
	if !ok {
		return nil, errors.New("secret not found")
	}
	secretData := st.(secret.Secret)
	return &secretData, nil
}

func (s *SecretRepositoryImpl) AddSecret(
	ctx context.Context,
	params secret.AddSecretParams,
) error {
	newSecret := secret.Secret{
		MasterKey:  params.MasterKey,
		Data:       params.Data,
		UnlockedAt: time.Now(),
	}
	s.cache.Set(params.StashID.String(), newSecret)
	return nil
}

func (s *SecretRepositoryImpl) RemoveSecret(
	ctx context.Context,
	params secret.RemoveSecretParams,
) error {
	s.cache.Delete(params.StashID.String())
	return nil
}
