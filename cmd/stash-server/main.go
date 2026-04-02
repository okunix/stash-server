package main

import (
	"github.com/okunix/stash-server/cmd/stash-server/cli"
	_ "github.com/okunix/stash-server/docs"
)

// @title						Stash API Server
// @version					0.0
// @description				Stash Password Manager API Server
// @basePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Bearer token. Format: "Bearer <token>"
func main() {
	cli.Execute()
}
