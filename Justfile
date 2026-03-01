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
        --build.cmd "go build -o {{out}} ./cmd/stash-server/main.go" \
        --build.bin "{{out}}" \
        --build.args_bin "--config=./server.yaml" \
        --build.exclude_dir "bin" \
        --misc.clean_on_exit "true"

compose:
    docker compose up --build --force-recreate

compose-down:
    docker compose down
