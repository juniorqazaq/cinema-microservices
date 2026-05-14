#!/usr/bin/env bash
set -euo pipefail
apt-get update -qq
apt-get install -y -qq protobuf-compiler >/dev/null
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
export PATH="$PATH:$(go env GOPATH)/bin"
mkdir -p gen/go
protoc \
  --proto_path=proto \
  --proto_path=third_party/googleapis \
  --go_out=gen/go \
  --go_opt=paths=source_relative \
  --go-grpc_out=gen/go \
  --go-grpc_opt=paths=source_relative \
  proto/movie/movie.proto
ls -la gen/go/movie
