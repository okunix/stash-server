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
