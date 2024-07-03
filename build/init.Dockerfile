FROM golang:1.22.5-alpine3.20 as builder
WORKDIR /build
COPY ./ ./
RUN <<EOF
export CGO_ENABLED=0
export GOOS=linux
go mod download
go build -C init
EOF

from alpine:latest as migrate
WORKDIR /app
ENV VERSION=v4.17.1 OS=linux ARCH=amd64
RUN <<EOF
apk add curl
curl -L https://github.com/golang-migrate/migrate/releases/download/${VERSION}/migrate.${OS}-${ARCH}.tar.gz | tar xvz
EOF

from alpine:latest
WORKDIR /app
RUN apk add yq-go postgresql-client
copy --from=builder /build/init/init .
copy --from=migrate /app/migrate /usr/bin
copy migration_files migration_files


ENTRYPOINT [ "./init" ]