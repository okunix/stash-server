package stash

import (
	"context"

	"github.com/google/uuid"
)

type UpdateStashParams struct {
	StashID     uuid.UUID
	Name        string
	Description *string
}

type CreateStashParams struct {
	Name          string
	Description   *string
	MaintainerID  uuid.UUID
	MasterKeyHash string
	MasterKeySalt string
}

type ListStashesParams struct {
	Limit        uint
	Offset       uint
	Search       string
	MaintainerID uuid.UUID
}

type CommitDataParams struct {
	StashID uuid.UUID
	Data    []byte
}

type Repository interface {
	ListStashes(ctx context.Context, params ListStashesParams) ([]*Stash, int64, error)
	GetStash(ctx context.Context, stashID uuid.UUID) (*Stash, error)
	CreateStash(ctx context.Context, params CreateStashParams) (*Stash, error)
	UpdateStash(ctx context.Context, params UpdateStashParams) (*Stash, error)
	DeleteStash(ctx context.Context, stashID uuid.UUID) error
	AddMember(ctx context.Context, userID uuid.UUID) error
	CommitData(ctx context.Context, params CommitDataParams) error
}
