package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io/fs"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

var sqliteConn *sql.DB

func SQLite() *sql.DB {
	return sqliteConn
}

func generate_uuid(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	return driver.String.ConvertValue(uuid.NewString())
}

func InitSQLite(ctx context.Context, path string, migrations fs.FS) error {
	sqlite.MustRegisterScalarFunction("generate_uuid", 0, generate_uuid)

	var err error
	sqliteConn, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	if err := sqliteConn.PingContext(ctx); err != nil {
		return err
	}
	pragmaStmt := `
		PRAGMA foreign_keys = 1;
	`
	_, err = sqliteConn.ExecContext(ctx, pragmaStmt)
	if err != nil {
		return err
	}
	return migrate(ctx, goose.DialectSQLite3, sqliteConn, migrations)
}

func migrate(ctx context.Context, dialect goose.Dialect, db *sql.DB, migrations fs.FS) error {
	provider, err := goose.NewProvider(dialect, db, migrations)
	if err != nil {
		return err
	}
	_, err = provider.Up(ctx)
	if err != nil {
		return err
	}
	slog.Info("database migration completed")
	return nil
}
