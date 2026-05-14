# Cinema microservices

Monorepository for a cinema booking platform: a **gRPC user and auth backend** (Go) and a **React web client** (Vite + TypeScript). The backend handles accounts, JWT sessions, and profiles; the frontend provides browsing, booking flows, and admin screens wired to your HTTP APIs.

## Repository layout

| Component | Path | Description |
| --------- | ---- | ----------- |
| **User service** | [`user-service/`](user-service/) | Go 1.22 gRPC API — register/login, JWT access & refresh, Redis-backed tokens, PostgreSQL persistence, health checks. Default gRPC port **50051**. |
| **Web frontend** | [`frontend/`](frontend/) | React 18 SPA — catalog, sessions, booking UI, profile, admin. Dev server via Vite (default **5173**). |
| **HTTP contract (gateway)** | [`GATEWAY_HTTP_CONTRACT.md`](GATEWAY_HTTP_CONTRACT.md) | Paths and gRPC mapping expected by the SPA; for the team member implementing the API gateway. |

Service-specific setup, environment variables, migrations, and `grpcurl` examples live in each package README.

## Prerequisites

- **Backend:** Go 1.22+, Docker & Docker Compose (recommended), or local PostgreSQL + Redis. Optional: `buf`, `protoc`, `grpcurl`.
- **Frontend:** Node.js **20+** (LTS recommended) and npm.

## Quick start

### User service (API)

From the service directory:

```bash
cd user-service
cp .env.example .env
# Edit .env — set JWT_SECRET, DB_*, REDIS_* as needed.

docker compose up -d --build
```

Or build and test without containers:

```bash
make build
make test
```

See **[user-service/README.md](user-service/README.md)** for migrations, protobuf generation, and gRPC usage.

### Web frontend

```bash
cd frontend
cp .env.example .env
# Set VITE_API_URL to your backend HTTP base URL (and optional TMDB keys for imagery).

npm install
npm run dev
```

Build and preview:

```bash
npm run build
npm run preview
```

Optional catalog tooling: `npm run generate:catalog` / `npm run scrape:catalog-paths` (see `frontend/scripts/`).

## Continuous integration

GitHub Actions (`.github/workflows/ci.yml`) runs on pushes and pull requests to `main`, `master`, and `sanat`:

- **user-service:** `go test` (with race), coverage, `gofmt` / `goimports`, `golangci-lint`, Docker image build.

Frontend CI is not wired in this workflow yet; run `npm run lint` and `npm run build` locally before opening a PR.

## Security

- Do **not** commit `.env` files or real secrets. Use `.env.example` as a template only.
- Rotate `JWT_SECRET` and database credentials for any shared or production environment.

## Contributing

1. Open a branch from the default target branch used in your fork/team.
2. Keep changes scoped; run tests and linters for the areas you touch (`make test` in `user-service`, `npm run lint` / `npm run build` in `frontend`).
3. Open a pull request with a short description of behavior and any new env vars or migration steps.

---

For deep dives, use the linked READMEs under [`user-service/`](user-service/) and [`frontend/`](frontend/).
