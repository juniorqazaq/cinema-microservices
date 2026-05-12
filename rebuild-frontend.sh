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

commit "chore(frontend): add package.json" \
  frontend/package.json

commit "chore(frontend): add npm lockfile" \
  frontend/package-lock.json

commit "chore(frontend): add TypeScript configs" \
  frontend/tsconfig.json \
  frontend/tsconfig.app.json \
  frontend/tsconfig.node.json

commit "chore(frontend): add Vite config and HTML shell" \
  frontend/vite.config.ts \
  frontend/index.html

commit "chore(frontend): add ESLint config" \
  frontend/eslint.config.js

commit "chore(frontend): add PostCSS and Tailwind" \
  frontend/postcss.config.js \
  frontend/tailwind.config.ts

commit "chore(frontend): add ignore rules and env template" \
  frontend/.gitignore \
  frontend/.env.example

commit "docs(frontend): add README" \
  frontend/README.md

commit "chore(frontend): add favicon asset" \
  frontend/public/favicon.svg

commit "chore(frontend): add public icons sprite" \
  frontend/public/icons.svg

commit "feat(frontend): add shared TypeScript types" \
  frontend/src/types/index.ts \
  frontend/src/types/movie.ts

commit "feat(frontend): add format helpers" \
  frontend/src/utils/format.ts

commit "feat(frontend): add JWT helpers" \
  frontend/src/utils/jwt.ts

commit "feat(frontend): add error handler utility" \
  frontend/src/utils/errorHandler.ts

commit "feat(frontend): add movie presentation helpers" \
  frontend/src/utils/moviePresentation.ts

commit "feat(frontend): add genre constants" \
  frontend/src/constants/genres.ts

commit "feat(frontend): add axios client" \
  frontend/src/api/axios.ts

commit "feat(frontend): add auth API" \
  frontend/src/api/auth.ts

commit "feat(frontend): add TMDB API helpers" \
  frontend/src/api/tmdb.ts

commit "feat(frontend): add movies API" \
  frontend/src/api/movies.ts

commit "feat(frontend): add bookings API" \
  frontend/src/api/bookings.ts

commit "feat(frontend): add admin API" \
  frontend/src/api/admin.ts

commit "feat(frontend): add auth store" \
  frontend/src/store/authStore.ts

commit "feat(frontend): add booking store" \
  frontend/src/store/bookingStore.ts

commit "feat(frontend): add movie image URL helpers" \
  frontend/src/lib/movieImages.ts

commit "feat(frontend): add catalog query filters" \
  frontend/src/lib/movieFilters.ts

commit "feat(frontend): add catalog adapter" \
  frontend/src/lib/catalogAdapter.ts

commit "feat(frontend): add movie catalog service" \
  frontend/src/lib/movieService.ts

commit "data(frontend): add now playing catalog slice" \
  frontend/src/data/nowPlaying.json

commit "data(frontend): add top rated catalog slice" \
  frontend/src/data/topRated.json

commit "data(frontend): add trending catalog slice" \
  frontend/src/data/trending.json

commit "data(frontend): add upcoming catalog slice" \
  frontend/src/data/upcoming.json

commit "chore(frontend): add TMDB path catalog script" \
  frontend/scripts/catalog-tmdb-paths.mjs

commit "chore(frontend): add catalog generator script" \
  frontend/scripts/generate-catalog.mjs

commit "chore(frontend): add TMDB scrape script" \
  frontend/scripts/scrape-tmdb-catalog-paths.mjs

commit "feat(frontend): add Badge component" \
  frontend/src/components/ui/Badge.tsx

commit "feat(frontend): add Button component" \
  frontend/src/components/ui/Button.tsx

commit "feat(frontend): add ErrorBanner component" \
  frontend/src/components/ui/ErrorBanner.tsx

commit "feat(frontend): add Input component" \
  frontend/src/components/ui/Input.tsx

commit "feat(frontend): add SectionTitle component" \
  frontend/src/components/ui/SectionTitle.tsx

commit "feat(frontend): add Spinner component" \
  frontend/src/components/ui/Spinner.tsx

commit "feat(frontend): add MovieHouse logo" \
  frontend/src/components/brand/MovieHouseLogo.tsx

commit "feat(frontend): add useAuth hook" \
  frontend/src/hooks/useAuth.ts

commit "feat(frontend): add useBooking hook" \
  frontend/src/hooks/useBooking.ts

commit "feat(frontend): add useMovies hook" \
  frontend/src/hooks/useMovies.ts

commit "feat(frontend): add useSessions hook" \
  frontend/src/hooks/useSessions.ts

commit "feat(frontend): add useTmdbMovie hook" \
  frontend/src/hooks/useTmdbMovie.ts

commit "feat(frontend): add Navbar" \
  frontend/src/components/layout/Navbar.tsx

commit "feat(frontend): add AppShell layout" \
  frontend/src/components/layout/AppShell.tsx

commit "feat(frontend): add protected route guard" \
  frontend/src/components/layout/ProtectedRoute.tsx

commit "feat(frontend): add admin route guard" \
  frontend/src/components/layout/AdminRoute.tsx

commit "feat(frontend): add HomeHero" \
  frontend/src/components/movie/HomeHero.tsx

commit "feat(frontend): add MovieCard" \
  frontend/src/components/movie/MovieCard.tsx

commit "feat(frontend): add MovieGrid" \
  frontend/src/components/movie/MovieGrid.tsx

commit "feat(frontend): add SessionPicker" \
  frontend/src/components/movie/SessionPicker.tsx

commit "feat(frontend): add CatalogMovieDetail" \
  frontend/src/components/movie/CatalogMovieDetail.tsx

commit "feat(frontend): add HomePage" \
  frontend/src/pages/movies/HomePage.tsx

commit "feat(frontend): add MovieDetailPage" \
  frontend/src/pages/movies/MovieDetailPage.tsx

commit "chore(frontend): add static image assets" \
  frontend/src/assets

commit "feat(frontend): add login page" \
  frontend/src/pages/auth/LoginPage.tsx

commit "feat(frontend): add register page" \
  frontend/src/pages/auth/RegisterPage.tsx

commit "feat(frontend): add seat map" \
  frontend/src/components/booking/SeatMap.tsx

commit "feat(frontend): add payment panel" \
  frontend/src/components/booking/PaymentPanel.tsx

commit "feat(frontend): add order summary" \
  frontend/src/components/booking/OrderSummary.tsx

commit "feat(frontend): add booking page" \
  frontend/src/pages/booking/BookingPage.tsx

commit "feat(frontend): add booking success page" \
  frontend/src/pages/booking/BookingSuccessPage.tsx

commit "feat(frontend): add ticket card" \
  frontend/src/components/profile/TicketCard.tsx

commit "feat(frontend): add profile page" \
  frontend/src/pages/profile/ProfilePage.tsx

commit "feat(frontend): add admin layout" \
  frontend/src/components/admin/AdminLayout.tsx

commit "feat(frontend): add hall form" \
  frontend/src/components/admin/HallForm.tsx

commit "feat(frontend): add movie form" \
  frontend/src/components/admin/MovieForm.tsx

commit "feat(frontend): add session form" \
  frontend/src/components/admin/SessionForm.tsx

commit "feat(frontend): add admin user table" \
  frontend/src/components/admin/UserTable.tsx

commit "feat(frontend): add admin dashboard page" \
  frontend/src/pages/admin/AdminDashboardPage.tsx

commit "feat(frontend): add admin movies page" \
  frontend/src/pages/admin/AdminMoviesPage.tsx

commit "feat(frontend): add admin users page" \
  frontend/src/pages/admin/AdminUsersPage.tsx

commit "feat(frontend): add admin bookings page" \
  frontend/src/pages/admin/AdminBookingsPage.tsx

commit "feat(frontend): add global styles" \
  frontend/src/index.css

commit "chore(frontend): add Vite client types" \
  frontend/src/vite-env.d.ts

commit "feat(frontend): add application routes" \
  frontend/src/App.tsx

commit "feat(frontend): add main entry" \
  frontend/src/main.tsx

git branch -D "$CURRENT_BRANCH" 2>/dev/null || true
git branch -m "$CURRENT_BRANCH"

echo "Done!"
git log --oneline --decorate -20
