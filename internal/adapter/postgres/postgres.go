package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var pgConn *sql.DB

func Postgres() *sql.DB {
	return pgConn
}

type PostgresInitParams struct {
	Migrations fs.FS
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	SSLMode    string
}

func (p PostgresInitParams) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
		p.SSLMode,
	)
}

func Init(ctx context.Context, params PostgresInitParams) error {
	var err error
	pgConn, err = sql.Open("postgres", params.String())
	if err != nil {
		return err
	}
	if err := pgConn.PingContext(ctx); err != nil {
		return err
	}
	return migrate(ctx, goose.DialectPostgres, pgConn, params.Migrations)
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
