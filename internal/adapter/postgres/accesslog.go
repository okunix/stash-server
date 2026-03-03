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

func scanAccessLogSQLModel(row scannable) (*accesslog.AccessLog, error) {
	var al accesslog.AccessLog
	err := row.Scan(&al.ID, &al.UserID, &al.StashID, &al.SecretName, &al.Action, &al.CreatedAt)
	return &al, err
}

const addLogStmt = `
INSERT INTO access_log (user_id, stash_id, secret_name, action) VALUES ($1, $2, $3, $4);
`

func (a *accessLogRepository) AddLog(
	ctx context.Context,
	params accesslog.CreateLogEntryParams,
) error {
	_, err := a.db.ExecContext(
		ctx,
		addLogStmt,
		params.UserID,
		params.StashID,
		params.SecretName,
		params.Action,
	)
	return err
}

func (a *accessLogRepository) RemoveLogs(
	ctx context.Context,
	timestamp *time.Time,
) (int64, error) {
	panic("unimplemented")
}

func (a *accessLogRepository) ListLogs(
	ctx context.Context,
	params accesslog.ListLogsParams,
) ([]*accesslog.AccessLog, int64, error) {
	panic("unimplemented")
}
