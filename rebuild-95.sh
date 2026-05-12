#!/bin/bash

set -e

CURRENT_BRANCH=$(git branch --show-current)

git checkout --orphan final-history
git reset

commit() {
  MSG=$1
  shift
  git add "$@"
  git commit -m "$MSG"
}

# ======================
# ROOT
# ======================

commit "chore: init monorepo structure" \
  README.md

# ======================
# USER SERVICE
# ======================

bash rebuild-user-service.sh

# ======================
# FRONTEND
# ======================

commit "chore(frontend): package setup" \
  frontend/package.json \
  frontend/package-lock.json

commit "chore(frontend): ts + vite config" \
  frontend/tsconfig.app.json \
  frontend/tsconfig.json \
  frontend/tsconfig.node.json \
  frontend/vite.config.ts

commit "chore(frontend): styling setup" \
  frontend/postcss.config.js \
  frontend/tailwind.config.ts \
  frontend/src/index.css

commit "chore(frontend): eslint + env" \
  frontend/eslint.config.js \
  frontend/.env.example \
  frontend/.gitignore

commit "docs(frontend): add README" \
  frontend/README.md

commit "feat(frontend): public assets" \
  frontend/public \
  frontend/src/assets

commit "feat(frontend): shared types" \
  frontend/src/types

commit "feat(frontend): utility helpers" \
  frontend/src/utils

commit "feat(frontend): constants" \
  frontend/src/constants

commit "feat(frontend): api layer" \
  frontend/src/api

commit "feat(frontend): state management" \
  frontend/src/store

commit "feat(frontend): lib layer" \
  frontend/src/lib

commit "data(frontend): catalogs" \
  frontend/src/data

commit "chore(frontend): scripts" \
  frontend/scripts

commit "feat(frontend): ui components" \
  frontend/src/components/ui

commit "feat(frontend): brand components" \
  frontend/src/components/brand

commit "feat(frontend): hooks" \
  frontend/src/hooks

commit "feat(frontend): layout system" \
  frontend/src/components/layout

commit "feat(frontend): movie components" \
  frontend/src/components/movie

commit "feat(frontend): auth pages" \
  frontend/src/pages/auth

commit "feat(frontend): movie pages" \
  frontend/src/pages/movies

commit "feat(frontend): booking system" \
  frontend/src/components/booking \
  frontend/src/pages/booking

commit "feat(frontend): profile module" \
  frontend/src/components/profile \
  frontend/src/pages/profile

commit "feat(frontend): admin dashboard" \
  frontend/src/components/admin \
  frontend/src/pages/admin

commit "feat(frontend): app bootstrap" \
  frontend/src/App.tsx \
  frontend/src/main.tsx \
  frontend/src/vite-env.d.ts

# ======================
# FINAL
# ======================

commit "chore: finalize cinema platform" .

git branch -D "$CURRENT_BRANCH" 2>/dev/null || true
git branch -m "$CURRENT_BRANCH"

echo "DONE"
git log --oneline --decorate -100
