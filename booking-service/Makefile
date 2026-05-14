PROTO_DIR     = proto
BOOKING_PROTO = $(PROTO_DIR)/booking/booking.proto

# Читаем .env для получения параметров БД
ifneq (,$(wildcard .env))
  include .env
  export
endif

# Строка подключения для migrate
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ─────────────────────────────────────────
# Proto
# ─────────────────────────────────────────

.PHONY: proto
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(BOOKING_PROTO)

.PHONY: proto-deps
proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# ─────────────────────────────────────────
# Миграции
# ─────────────────────────────────────────

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

.PHONY: migrate-drop
migrate-drop:
	migrate -path migrations -database "$(DB_URL)" drop -f

# ─────────────────────────────────────────
# Установка migrate (один раз)
# ─────────────────────────────────────────

.PHONY: migrate-install
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ─────────────────────────────────────────
# Сборка и запуск
# ─────────────────────────────────────────

.PHONY: build
build:
	go build -o bin/booking-service ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go
