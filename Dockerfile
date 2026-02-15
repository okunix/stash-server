FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app
COPY . .
RUN apk add --no-cache -X https://dl-cdn.alpinelinux.org/alpine/edge/community just 
RUN just build

FROM alpine:3.23
COPY --from=builder /app/bin/* .
HEALTHCHECK CMD [ "sh", "-c" "wget -qO- http://localhost/health >/dev/null || exit 1"]

CMD [ "/stash-server" ]
