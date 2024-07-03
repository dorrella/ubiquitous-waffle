FROM golang:latest as builder
WORKDIR /build
COPY ./ ./
RUN <<EOF
export CGO_ENABLED=0
export GOOS=linux
go mod download
go build -C service
EOF

from alpine:latest
WORKDIR /app
RUN apk add yq-go postgresql-client
copy --from=builder /build/service/service .

ENTRYPOINT [ "./service" ]