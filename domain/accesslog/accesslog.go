package accesslog

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/domain/user"
)

type AccessLog struct {
	ID         uuid.UUID   `json:"id"`
	User       user.User   `json:"user"`
	Stash      stash.Stash `json:"stash"`
	SecretName string      `json:"secret_name"`
	Action     string      `json:"action"`
	CreatedAt  time.Time   `json:"created_at"`
}

type ListLogsParams struct {
	StashID uuid.UUID
	Limit   uint
	Offset  uint
}

type DeleteLogsParams struct {
	Timestamp *time.Time
}

type CreateLogEntryParams struct {
	UserID     uuid.UUID
	StashID    uuid.UUID
	SecretName string
	Action     string
}

type Repository interface {
	ListLogs(ctx context.Context, params ListLogsParams) ([]*AccessLog, error)
	DeleteLogs(ctx context.Context, params DeleteLogsParams) (int64, error)
	CreateLogEntry(ctx context.Context, params CreateLogEntryParams) error
}
