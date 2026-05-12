#!/usr/bin/env bash
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

CURRENT_BRANCH="$(git branch --show-current)"

git checkout --orphan clean-history
git reset

commit() {
  local message="$1"
  shift

  git add "$@"
  git commit -m "$message"
}

# =========================
# ROOT
# =========================
commit "chore: init monorepo structure (frontend + user-service)" \
  README.md \
  .gitignore

# =========================
# USER SERVICE (backend - 14 commits condensed)
# =========================
commit "chore(user-service): bootstrap backend service" \
  user-service/go.mod \
  user-service/go.sum \
  user-service/Makefile \
  user-service/.golangci.yml \
  user-service/README.md \
  user-service/.env.example

commit "feat(user-service): gRPC API + proto definitions" \
  user-service/buf.yaml \
  user-service/buf.gen.yaml \
  user-service/proto \
  user-service/third_party

commit "feat(user-service): generated protobuf code" \
  user-service/gen/go/user

commit "feat(user-service): domain layer (entities + repos)" \
  user-service/internal/domain \
  user-service/internal/repository/mocks

commit "feat(user-service): config and environment loader" \
  user-service/internal/config

commit "feat(user-service): database layer (postgres + migrations)" \
  user-service/migrations \
  user-service/internal/repository/postgres \
  user-service/internal/repository/redis

commit "feat(user-service): security utilities (auth, validation)" \
  user-service/internal/security \
  user-service/internal/validation

commit "feat(user-service): authentication and profile usecases" \
  user-service/internal/usecase \
  user-service/tests/usecase

commit "feat(user-service): gRPC delivery layer" \
  user-service/internal/delivery/grpc

commit "feat(user-service): application bootstrap + health checks" \
  user-service/internal/health \
  user-service/cmd

commit "chore(user-service): devops (docker + ci)" \
  user-service/Dockerfile \
  user-service/docker-compose.yml \
  .github/workflows/ci.yml \
  user-service/tests/repository

# =========================
# FRONTEND (81 commits condensed)
# =========================
commit "feat(frontend): full application (UI + API + state + routing)" \
  frontend

# =========================
# FINAL SNAPSHOT
# =========================
commit "chore: finalize cinema platform (frontend + user-service)" \
  .

git branch -D "$CURRENT_BRANCH" 2>/dev/null || true
git branch -m "$CURRENT_BRANCH"

echo "Done!"
git log --oneline --decorate -20
