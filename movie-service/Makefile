BINARY      := movie-service
CMD_DIR     := ./cmd
PROTO_DIR   := proto
GO_OUT      := gen/go
MIGRATE_BIN := migrate
DB_URL      ?= postgres://cinema:cinema_secret@localhost:5432/movie_db?sslmode=disable
IMAGE       ?= movie-service:latest

.PHONY: build
build:
	go build -ldflags="-w -s" -o bin/$(BINARY) $(CMD_DIR)/main.go

.PHONY: run
run:
	go run $(CMD_DIR)/main.go

.PHONY: proto
proto:
	@mkdir -p $(GO_OUT)
	PATH="$$(pwd)/.bin:$$PATH" protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=third_party/googleapis \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/movie/movie.proto

.PHONY: test
test:
	go test ./... -count=1

.PHONY: test-cover
test-cover:
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: test-unit
test-unit:
	go test ./tests/usecase/... -count=1 -v

.PHONY: test-integration
test-integration:
	go test -tags=integration ./tests/repository/... -count=1 -v -timeout 120s

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_BIN) -path ./migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -path ./migrations -database "$(DB_URL)" down 1

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	gofmt -w .
	goimports -w .

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) .

.PHONY: install-tools
install-tools:
	mkdir -p .bin
	GOBIN=$$(pwd)/.bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
	GOBIN=$$(pwd)/.bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/golang/mock/mockgen@v1.6.0
	go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

.PHONY: generate-mocks
generate-mocks:
	@echo "Use mockgen with -source=internal/domain/movie.go (or install mockgen); tests use tests/testutil fakes."

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: help
help:
	@echo "Targets: build run proto test migrate-up migrate-down docker-build install-tools generate-mocks tidy"
