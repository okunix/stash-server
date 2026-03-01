package postgres

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/accesslog"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type accessLogRepository struct {
	db *sql.DB
}

func NewAccessLogRepository(db *sql.DB) ports.AccessLogRepository {
	return &accessLogRepository{db: db}
}

func (a *accessLogRepository) AddLog(
	ctx context.Context,
	params accesslog.CreateLogEntryParams,
) error {
	panic("unimplemented")
}

func (a *accessLogRepository) ListLogs(
	ctx context.Context,
	params accesslog.ListLogsParams,
) ([]*accesslog.AccessLog, int64, error) {
	panic("unimplemented")
}

func (a *accessLogRepository) RemoveLogs(
	ctx context.Context,
	timestamp *time.Time,
) (int64, error) {
	panic("unimplemented")
}
