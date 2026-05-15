# Cinema microservices

Monorepository for a cinema booking platform: **Go gRPC microservices** (users, catalog, bookings, notifications) and a **React** web client (Vite + TypeScript). The browser app talks to an **HTTP API** (typically implemented as a separate API gateway); it does not call gRPC directly. Set `VITE_API_URL` in `frontend/.env` to that HTTP base URL (default in code: `http://localhost:8080`).

---

## Repository layout

| Component | Path | Default port / transport |
| --------- | ---- | ------------------------- |
| **User service** | [`user-service/`](user-service/) | gRPC **50051** |
| **Movie service** | [`movie-service/`](movie-service/) | gRPC **50052** |
| **Booking service** | [`booking-service/`](booking-service/) | gRPC **50053** |
| **API gateway** | [`api-gateway/`](api-gateway/) | HTTP **8080** |
| **Notification service** | [`notification-service/`](notification-service/) | NATS subscriber + SMTP (no public HTTP) |
| **Web frontend** | [`frontend/`](frontend/) | Vite dev **5173** |

---

## User service

**Path:** [`user-service/`](user-service/) · **Details:** [user-service/README.md](user-service/README.md)

- **Role:** Accounts, authentication (JWT access + refresh), profiles, token blacklist, admin user operations.
- **Stack:** Go 1.22+, PostgreSQL, Redis, Clean Architecture, gRPC.
- **Proto:** [`user-service/proto/user/user.proto`](user-service/proto/user/user.proto) — service `UserService` (register, login, logout, profile, `ValidateToken` for gateways, refresh, password change, admin list/ban, etc.).
- **Codegen:** `make proto` in `user-service/` → `gen/go/user/`.
- **Run:** `cp .env.example .env`, then `docker compose up` from `user-service/` or `make build` / `make test`.

---

## Movie service

**Path:** [`movie-service/`](movie-service/)

- **Role:** Movies, halls, seats, sessions; listing and admin-style mutations guarded by `requester_role` in proto where applicable.
- **Stack:** Go, PostgreSQL, Redis, gRPC.
- **Proto:** [`movie-service/proto/movie/movie.proto`](movie-service/proto/movie/movie.proto) — service `MovieService`.
- **Codegen:** follow `movie-service` Makefile / project conventions (`gen/go/movie/`).
- **Config:** `DB_*`, `REDIS_ADDR`, `GRPC_PORT` (default **50052**). See `movie-service/internal/config` and `.env.example` if present.

---

## Booking service

**Path:** [`booking-service/`](booking-service/)

- **Role:** Bookings, payments, seat checks, admin booking list/stats; publishes domain events to NATS for notifications.
- **Stack:** Go, PostgreSQL, NATS, gRPC.
- **Proto:** [`booking-service/proto/booking/booking.proto`](booking-service/proto/booking/booking.proto) — service `BookingService`.
- **Codegen:** `make proto` from `booking-service/`.
- **Run:** `DATABASE_URL`, `NATS_URL`, `GRPC_PORT` (default **50053**). Root [`docker-compose.yml`](docker-compose.yml) includes Postgres + NATS + booking for local integration.

---

## Notification service

**Path:** [`notification-service/`](notification-service/)

- **Role:** Consumes NATS events and sends email (SMTP) using HTML templates under `internal/templates/`.
- **Stack:** Go, NATS client, SMTP.
- **Subjects:** `booking.created`, `booking.cancelled`, `payment.confirmed` (payloads must include `email` when mail should be sent; empty `email` skips send).
- **Env:** `NATS_URL`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` (see `notification-service/cmd/main.go` and root `docker-compose.yml` for examples).

---

## Web frontend

**Path:** [`frontend/`](frontend/) · **Details:** [frontend/README.md](frontend/README.md)

- **Role:** Catalog, sessions, booking flow, profile, admin screens.
- **API client:** `frontend/src/api/` — Axios, auth refresh queue, resource modules. Expects HTTP JSON with a `data` envelope where noted in `axios.ts`.
- **Run:** `cd frontend && cp .env.example .env && npm install && npm run dev`.

---

## HTTP API and gateway

The HTTP gateway lives in [`api-gateway/`](api-gateway/). The SPA’s expected routes and payloads are defined by the TypeScript client:

- **Source of truth:** [`frontend/src/api/`](frontend/src/api/) (especially `auth.ts`, `movies.ts`, `bookings.ts`, `admin.ts`) and [`frontend/src/api/axios.ts`](frontend/src/api/axios.ts).

**Mapping:** gateway listens on HTTP 8080, validates `Authorization: Bearer …` via **UserService.ValidateToken**, then forwards to:

| gRPC target | Address (local defaults) | Proto |
| ----------- | ------------------------ | ----- |
| User | `localhost:50051` | `user.UserService` |
| Movie | `localhost:50052` | `movie.MovieService` |
| Booking | `localhost:50053` | `booking.BookingService` |

**Conventions used by the client:** success body `{ "data": … }`; errors `{ "error": "…" }`; register/login/refresh do not require Bearer; protected calls attach Bearer and rely on `/auth/refresh` on 401 (see `axios.ts`). Booking gRPC calls pass **`user_email`** from the validated user when creating/cancelling bookings and confirming payments so notification-service can email users.

**Gateway run:**

```bash
cd api-gateway
go run ./cmd/api-gateway
```

The gateway uses in-memory rate limits by IP and authenticated user. Defaults are documented in [`api-gateway/README.md`](api-gateway/README.md).

---

## Prerequisites

- **Backend:** Go 1.22+, Docker & Docker Compose (optional), PostgreSQL, Redis where a service needs it, NATS for booking/notification. Optional: `protoc`, `grpcurl`.
- **Frontend:** Node.js **20+**, npm.

---

## Quick start (local)

**User service**

```bash
cd user-service
cp .env.example .env
docker compose up -d --build
```

**Frontend**

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

**Booking stack (Postgres + NATS + booking + notification)** from repo root:

```bash
docker compose up -d --build
```

---

## Continuous integration

GitHub Actions (`.github/workflows/ci.yml`) runs on pushes and PRs to `main`, `master`, and `sanat` (user-service tests, lint, Docker build). Run `npm run lint` and `npm run build` in `frontend/` locally for frontend changes.

---

## Security

- Do **not** commit `.env` or real secrets; use `*.env.example` only as templates.
- Rotate `JWT_SECRET` and database credentials for shared or production environments.

---

## Contributing

1. Branch from your team’s default branch.
2. Scope changes; run tests and linters for touched areas (`make test` in `user-service`, `go test ./...` in other services, `npm run lint` / `npm run build` in `frontend`).
3. Open a PR describing behavior and any new environment variables or migration steps.
