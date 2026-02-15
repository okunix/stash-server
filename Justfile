set dotenv-load

_default:
    @just -f {{justfile()}} --list

migrations_dir := "./migrations"

goose_cmd := f"goose $DB_DRIVER $DB_URL -dir {{migrations_dir}} -v"

# create new migration
new-migration name:
    goose create -dir {{migrations_dir}} -s {{name}} sql

# apply migration
migrate-up version="":
    #!/usr/bin/env bash
    if [ -z {{version}} ]; then
        {{goose_cmd}} up
    else
        {{goose_cmd}} up-to {{version}}
    fi


# downgrade migration
migrate-down version="":
    #!/usr/bin/env bash
    if [ -z {{version}} ]; then
        {{goose_cmd}} down
    else
        {{goose_cmd}} down-to {{version}}
    fi

_go_deps:
    go mod tidy
    go mod download
    [ -x "$(command -v sqlc)" ] || \
        go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    sqlc generate

export GOOS:="linux"
export GOARCH:=`uname -m | sed -E 's/^x86_64$/amd64/g'`
export CGO_ENABLED:="0"

# build project
build out="./bin/stash-server": _go_deps
    go build -o {{out}} .

# run binary
run config="config.yaml": build
    ./bin/stash-server --config {{config}}

docker_image_name := "stash-server"

# build docker image
docker-build tag="latest":
    docker build -t {{docker_image_name}}:{{tag}} .
