package secret

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// in-memory data structure
type Secret struct {
	MasterKey  string            `json:"master_key"`
	Data       map[string]string `json:"data"`
	UnlockedAt time.Time         `json:"unlocked_at"`
}

type AddSecretParams struct {
	StashID      uuid.UUID
	MaintainerID uuid.UUID
	MasterKey    string
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
