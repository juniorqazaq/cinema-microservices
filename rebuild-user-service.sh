#!/usr/bin/env bash
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

CURRENT_BRANCH="$(git branch --show-current)"

git checkout --orphan temp-history
git reset

commit() {
  local message="$1"
  shift

  git add "$@"
  git commit -m "$message"
}

commit "chore(user-service): bootstrap project" \
  user-service/go.mod \
  user-service/go.sum \
  user-service/Makefile \
  user-service/.golangci.yml \
  user-service/README.md \
  user-service/.env.example

commit "feat(proto): add gRPC API and buf config" \
  user-service/buf.yaml \
  user-service/buf.gen.yaml \
  user-service/proto \
  user-service/third_party

commit "chore(proto): generate protobuf code" \
  user-service/gen/go/user

commit "feat(domain): add core user entity and repository contracts" \
  user-service/internal/domain/user.go \
  user-service/internal/repository/mocks

commit "feat(config): add environment configuration loader" \
  user-service/internal/config/config.go

commit "feat(db): add PostgreSQL migrations" \
  user-service/migrations

commit "feat(repository): implement PostgreSQL repository" \
  user-service/internal/repository/postgres

commit "feat(repository): add Redis token repository" \
  user-service/internal/repository/redis

commit "feat(security): add password hashing and validation" \
  user-service/internal/security \
  user-service/internal/validation

commit "feat(auth): implement authentication usecases" \
  user-service/internal/usecase/auth.go \
  user-service/tests/usecase/auth_test.go

commit "feat(profile): implement profile management usecases" \
  user-service/internal/usecase/profile.go \
  user-service/tests/usecase/profile_test.go

commit "feat(grpc): add delivery layer and handlers" \
  user-service/internal/delivery/grpc

commit "feat(app): wire dependencies and health checks" \
  user-service/internal/health \
  user-service/cmd/main.go

commit "chore(devops): add Docker and CI pipeline" \
  user-service/Dockerfile \
  user-service/docker-compose.yml \
  .github/workflows/ci.yml \
  user-service/tests/repository/postgres_test.go

git branch -D "$CURRENT_BRANCH" 2>/dev/null || true
git branch -m "$CURRENT_BRANCH"

echo "Done!"
git log --oneline --decorate -15
