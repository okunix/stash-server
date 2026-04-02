set dotenv-load

migrations_dir := "./migrations"

docker_image_name := "stash-server"

export GOOS:="linux"
export GOARCH:=`uname -m | sed -E 's/^x86_64$/amd64/g'`
export CGO_ENABLED:="0"

_default:
    @just -f {{justfile()}} --list

_go_deps:
    go mod tidy
    go mod download
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g ./cmd/stash-server/main.go

# create new migration
new-migration name:
    goose create -dir {{migrations_dir}} -s {{name}} sql

# build project
build out="./bin/stash-server": _go_deps
    go build -o {{out}} ./cmd/stash-server/main.go

# run binary
run config="server.yaml": build
    ./bin/stash-server --config {{config}}

# build docker image
docker-build tag="latest":
    docker build -t {{docker_image_name}}:{{tag}} .

# run with air hot-reload
run-air out="./bin/stash-server":
    air --tmp_dir="./bin" \
        --build.cmd "just build {{out}}" \
        --build.bin "{{out}}" \
        --build.args_bin "--config=./server.yaml" \
        --build.exclude_dir "bin" \
        --build.exclude_dir "docs" \
        --build.include_file "server.yaml" \
        --misc.clean_on_exit "true"

# run docker compose with build and force-recreate options
compose:
    docker compose up --build --force-recreate

# run docker compose down
compose-down:
    docker compose down
