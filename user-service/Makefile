BINARY      := user-service
CMD_DIR     := ./cmd
PROTO_DIR   := proto
GO_OUT      := gen/go
MIGRATE_BIN := migrate
DB_URL      ?= postgres://cinema:cinema_secret@localhost:5432/cinema_db?sslmode=disable
IMAGE       ?= user-service:latest

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
		$(PROTO_DIR)/user/user.proto

.PHONY: test
test:
	go test ./... -race -count=1

.PHONY: test-cover
test-cover:
	go test ./... -race -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: test-unit
test-unit:
	go test ./tests/usecase/... -race -count=1 -v

.PHONY: test-integration
test-integration:
	go test -tags=integration ./tests/repository/... -race -count=1 -v -timeout 120s

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_BIN) -path ./migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -path ./migrations -database "$(DB_URL)" down 1

.PHONY: migrate-drop
migrate-drop:
	$(MIGRATE_BIN) -path ./migrations -database "$(DB_URL)" drop -f

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

.PHONY: docker-up
docker-up:
	docker compose up -d --build

.PHONY: docker-down
docker-down:
	docker compose down -v

.PHONY: docker-logs
docker-logs:
	docker compose logs -f user-service

.PHONY: generate-mocks
generate-mocks:
	go run github.com/golang/mock/mockgen@v1.6.0 \
		-source=internal/domain/user.go \
		-destination=internal/repository/mocks/mock_repository.go \
		-package=mocks

.PHONY: install-tools
install-tools:
	mkdir -p .bin
	GOBIN=$$(pwd)/.bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
	GOBIN=$$(pwd)/.bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/golang/mock/mockgen@v1.6.0
	go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

.PHONY: download-googleapis
download-googleapis:
	mkdir -p third_party/googleapis/google/protobuf
	curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/timestamp.proto \
		-o third_party/googleapis/google/protobuf/timestamp.proto

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: help
help:
	@echo "Targets: build run proto test test-cover test-unit test-integration lint fmt migrate-up migrate-down docker-build docker-up docker-down install-tools download-googleapis tidy generate-mocks"
