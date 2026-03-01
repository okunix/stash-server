package accesslog

import (
	"time"

	"github.com/google/uuid"
)

type AccessLog struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user"`
	StashID    uuid.UUID `json:"stash"`
	SecretName string    `json:"secret_name"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListLogsParams struct {
	StashID uuid.UUID
	Limit   uint
	Offset  uint
}

type CreateLogEntryParams struct {
	UserID     uuid.UUID
	StashID    uuid.UUID
	SecretName string
	Action     string
}
