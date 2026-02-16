package stash

import (
	"context"
	"os/user"
	"time"

	"github.com/google/uuid"
)

type Stash struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"desc"`
	Maintainer    user.User   `json:"maintainer"`
	Members       []user.User `json:"members"`
	MasterKeyHash string      `json:"-"`
	MasterKeySalt string      `json:"-"`
	Content       []byte      `json:"content"`
	CreatedAt     time.Time   `json:"created_at"`
}

type AddMemberParams struct {
	UserID uuid.UUID
}

type DeleteStashParams struct {
	StashID uuid.UUID
}

type UpdateStashParams struct {
	ID          uuid.UUID
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
	ID uuid.UUID
}

type ListStashesParams struct {
	Limit  int64
	Offset int64
	Total  int64
	Search string
}

type CommitContentParams struct {
	Content []byte
}

type StashRepository interface {
	ListStashes(ctx context.Context, params ListStashesParams) ([]*Stash, error)
	GetStash(ctx context.Context, params GetStashParams) (*Stash, error)
	CreateStash(ctx context.Context, params CreateStashParams) (*Stash, error)
	UpdateStash(ctx context.Context, params UpdateStashParams) (*Stash, error)
	DeleteStash(ctx context.Context, params DeleteStashParams) error
	AddMember(ctx context.Context, params AddMemberParams) error
	CommitContent(ctx context.Context, params CommitContentParams) error
}
