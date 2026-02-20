package stash

import (
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
	Data          []byte      `json:"data"`
	CreatedAt     time.Time   `json:"created_at"`
}
