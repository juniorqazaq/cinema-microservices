# User Service (Cinema Booking)

Run all `make` / `docker compose` commands from this `user-service/` directory. Monorepo overview: [README](../README.md).

gRPC user and authentication service: registration, login, JWT access and refresh tokens, profiles, and admin operations. Stack: Go 1.22, PostgreSQL, Redis, Clean Architecture.

## Architecture

- **Domain** (`internal/domain`): entities, repository interfaces, use case contracts, sentinel errors.
- **Use cases** (`internal/usecase`): auth (JWT, bcrypt) and profile/admin workflows.
- **Repositories** (`internal/repository/postgres`, `redis`): PostgreSQL via pgx, Redis for refresh tokens, blacklist, and user cache.
- **Delivery** (`internal/delivery/grpc`): gRPC handlers, interceptors (logging, timing, recovery), reflection for grpcurl.
- **Cross-cutting** (`internal/errors`, `internal/validation`, `internal/security`, `internal/health`): gRPC error mapping, UUID and email rules, constant-time-friendly password checks, health checks.

```
cmd/main.go
internal/
  config/
  domain/
  delivery/grpc/
  errors/
  health/
  repository/
  security/
  usecase/
  validation/
proto/user/
gen/go/user/
migrations/
```

## Setup

1. Install Go 1.22+, Docker (optional), `protoc`, and plugins (`make install-tools`, `make download-googleapis`).
2. Copy `.env.example` to `.env` and set secrets (especially `JWT_SECRET`).
3. Start PostgreSQL and Redis (see Docker section) and run migrations: `make migrate-up`.
4. Generate code if you change protos: `make proto`.

## Environment

| Variable | Description |
|----------|-------------|
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | PostgreSQL |
| `REDIS_ADDR`, `REDIS_PASSWORD` | Redis |
| `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `JWT_ISSUER`, `JWT_AUDIENCE` | JWT |
| `GRPC_PORT` | gRPC listen port (default `50051`) |

## Docker

```bash
cp .env.example .env
docker compose up -d --build
```

Compose brings up Postgres, Redis, a one-shot migrate container, and `user-service`.

## Migrations

Uses [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
make migrate-up
make migrate-down
```

Default `DB_URL` for Make is in the `Makefile`; override with `DB_URL=... make migrate-up`.

## gRPC and grpcurl

List services:

```bash
grpcurl -plaintext localhost:50051 list
```

Health (checks DB and Redis):

```bash
grpcurl -plaintext -d '{}' localhost:50051 grpc.health.v1.Health/Check
```

Example login (adjust JSON):

```bash
grpcurl -plaintext -d '{"email":"user@example.com","password":"yourpassword"}' \
  localhost:50051 user.UserService/Login
```

Admin-only RPCs (`GetAllUsers`, `GetUserByEmail`, `BanUser`) require a valid `admin_id` of a user with role `ADMIN` in the request messages.

## Testing

```bash
make test              # unit tests (default ./...)
make test-cover        # coverage summary
make test-integration  # Postgres via testcontainers; needs Docker
```

## Development

```bash
make fmt
make lint
make docker-build
```

## Security notes

- Passwords are hashed with bcrypt (cost 12). Failed login performs a dummy bcrypt compare to reduce user enumeration timing signals.
- Emails are normalized (trim + lowercase) on register, login, and admin lookup.
- User IDs in mutating flows are validated as UUIDs.
- JWTs are validated for issuer, audience, expiry, and HMAC-SHA256 only.
- Structured logs avoid passwords, secrets, and raw tokens.
