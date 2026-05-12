# Cinema web frontend

Single-page app for the cinema microservices project: movie catalog, sessions, seat booking, payments (UI), user profile, and admin screens. Built with **React 18**, **TypeScript**, **Vite**, **Tailwind CSS**, **TanStack Query**, **Zustand**, and **React Router**.

Monorepo overview: [../README.md](../README.md).

## Requirements

- **Node.js 20+** (LTS recommended)
- npm (ships with Node)

## Setup

```bash
cp .env.example .env
npm install
```

### Environment variables

Copy [.env.example](.env.example) to `.env` and adjust:

| Variable | Description |
| -------- | ----------- |
| `VITE_API_URL` | Base URL of the HTTP API the app calls (default in code: `http://localhost:8080`). |
| `VITE_TMDB_READ_ACCESS_TOKEN` | Optional [TMDB](https://www.themoviedb.org/settings/api) read token for posters and hero imagery. |
| `VITE_TMDB_API_KEY` | Optional TMDB API key; the app prefers the Bearer token when both are set. |

Never commit `.env` or real secrets.

## Scripts

| Command | Description |
| ------- | ----------- |
| `npm run dev` | Start Vite dev server (HMR). |
| `npm run build` | Typecheck and production build to `dist/`. |
| `npm run preview` | Serve the production build locally. |
| `npm run lint` | Run ESLint on the project. |
| `npm run generate:catalog` | Regenerate local catalog JSON via `scripts/generate-catalog.mjs`. |
| `npm run scrape:catalog-paths` | Scrape TMDB path lists via `scripts/scrape-tmdb-catalog-paths.mjs`. |

## Project layout (short)

- `src/api/` — Axios client, auth refresh queue, resource modules.
- `src/pages/` — Route-level views (movies, auth, booking, profile, admin).
- `src/components/` — Layout, UI primitives, feature blocks.
- `src/hooks/` — React Query–backed data hooks.
- `src/store/` — Client state (auth, booking draft).
- `src/lib/` — Catalog helpers, filters, image URLs.
- `src/data/` — Bundled catalog JSON slices used for browsing when TMDB is optional.

## Production notes

- Point `VITE_API_URL` at your deployed API and rebuild; Vite inlines env at build time.
- Serve `dist/` behind HTTPS in production; configure CORS on the API for your origin.

## Troubleshooting

- **401 / repeated logout:** confirm the API URL, CORS, and that refresh-token endpoints match what `src/api/axios.ts` expects.
- **Missing posters:** set TMDB variables or rely on bundled catalog paths where available.
