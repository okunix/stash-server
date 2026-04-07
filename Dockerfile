FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app
COPY . .
RUN apk add --no-cache -X https://dl-cdn.alpinelinux.org/alpine/edge/community just 
RUN just build

FROM alpine:3.23
COPY --from=builder /app/bin/* .
COPY --from=builder /app/config.example.yml /etc/stash/server.yaml

HEALTHCHECK CMD [ "sh", "-c" "wget -qO- http://localhost:7878/health >/dev/null || exit 1"]
EXPOSE 7878

CMD [ "/stash-server" ]
