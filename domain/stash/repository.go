package stash

import (
	"context"

	"github.com/google/uuid"
)

type AddMemberParams struct {
	UserID uuid.UUID
}

type DeleteStashParams struct {
	StashID uuid.UUID
}

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

type GetStashParams struct {
	StashID uuid.UUID
}

type ListStashesParams struct {
	Limit        uint
	Offset       uint
	Total        uint
	Search       string
	MaintainerID uuid.UUID
}

type CommitDataParams struct {
	StashID uuid.UUID
	Data    []byte
}

type Repository interface {
	ListStashes(ctx context.Context, params ListStashesParams) ([]*Stash, error)
	GetStash(ctx context.Context, params GetStashParams) (*Stash, error)
	CreateStash(ctx context.Context, params CreateStashParams) (*Stash, error)
	UpdateStash(ctx context.Context, params UpdateStashParams) (*Stash, error)
	DeleteStash(ctx context.Context, params DeleteStashParams) error
	AddMember(ctx context.Context, params AddMemberParams) error
	CommitData(ctx context.Context, params CommitDataParams) error
}
