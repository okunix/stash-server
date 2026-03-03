package stash

import (
	"time"

	"github.com/google/uuid"
)

type Stash struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"desc"`
	MaintainerID  uuid.UUID `json:"maintainer_id"`
	MasterKeyHash string    `json:"-"`
	EncryptedData string    `json:"data"`
	CreatedAt     time.Time `json:"created_at"`
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

type ListStashesParams struct {
	Limit        uint
	Offset       uint
	Search       string
	MaintainerID uuid.UUID
}

type CommitDataParams struct {
	StashID uuid.UUID
	Data    string
}

type AddMemberParams struct {
	UserID  uuid.UUID
	StashID uuid.UUID
}
