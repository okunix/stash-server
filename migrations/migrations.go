package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var migrationsFS embed.FS

func Migrations() fs.FS {
	return migrationsFS
}
