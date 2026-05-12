import type { CatalogMovie } from '../types/movie'

export const TMDB_IMAGE_BASE = 'https://image.tmdb.org/t/p'

export type TmdbImageSize =
  | 'w200'
  | 'w300'
  | 'w342'
  | 'w500'
  | 'w780'
  | 'original'

const FALLBACK_POSTER =
  'data:image/svg+xml,' +
  encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="500" height="750" viewBox="0 0 500 750"><rect fill="#101827" width="500" height="750"/><text x="50%" y="50%" fill="#64748b" font-family="system-ui,sans-serif" font-size="22" text-anchor="middle" dominant-baseline="middle">No poster</text></svg>`,
  )

const FALLBACK_BACKDROP =
  'data:image/svg+xml,' +
  encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="1280" height="720" viewBox="0 0 1280 720"><rect fill="#050a18" width="1280" height="720"/><text x="50%" y="50%" fill="#334155" font-family="system-ui,sans-serif" font-size="28" text-anchor="middle" dominant-baseline="middle">No backdrop</text></svg>`,
  )

/** Non-TMDB absolute URLs (cast avatars, legacy data). */
export function isAbsoluteImageUrl(path: string): boolean {
  return /^https?:\/\//i.test(path.trim())
}

/**
 * Build a TMDB CDN URL from a **file path** (`/abc.jpg`) or pass through external `https:` URLs.
 * For TMDB paths, `size` is the `/t/p/{size}/` segment (w200 … original).
 */
export function tmdbImage(
  path: string | null | undefined,
  size: TmdbImageSize = 'w500',
): string | null {
  if (!path?.trim()) return null
  const raw = path.trim()
  if (isAbsoluteImageUrl(raw)) {
    if (raw.includes('image.tmdb.org')) {
      const m = raw.match(/\/t\/p\/([^/]+)(\/[^?]+)/)
      if (m) return `${TMDB_IMAGE_BASE}/${size}${m[2]}`
    }
    return raw
  }
  const file = raw.startsWith('/') ? raw : `/${raw}`
  return `${TMDB_IMAGE_BASE}/${size}${file}`
}

function widthToPosterSize(width: number): TmdbImageSize {
  if (width <= 200) return 'w200'
  if (width <= 300) return 'w300'
  if (width <= 342) return 'w342'
  if (width <= 500) return 'w500'
  return 'w780'
}

function widthToBackdropSize(width: number): TmdbImageSize {
  if (width <= 780) return 'w780'
  return 'original'
}

/**
 * Resolve catalog `poster_path` / `backdrop_path` with sensible TMDB size for layout width.
 * External non-TMDB URLs are returned unchanged.
 */
export function optimizeCatalogImageUrl(
  path: string | null | undefined,
  opts: { width: number; role: 'poster' | 'backdrop' },
): string {
  if (!path?.trim()) {
    return opts.role === 'poster' ? FALLBACK_POSTER : FALLBACK_BACKDROP
  }
  const raw = path.trim()
  if (isAbsoluteImageUrl(raw) && !raw.includes('image.tmdb.org')) {
    return raw
  }
  const size =
    opts.role === 'poster'
      ? widthToPosterSize(opts.width)
      : widthToBackdropSize(opts.width)
  return tmdbImage(raw, size) ?? (opts.role === 'poster' ? FALLBACK_POSTER : FALLBACK_BACKDROP)
}

export function posterUrl(
  movie: Pick<CatalogMovie, 'poster_path'>,
  widthOrSize: number | TmdbImageSize = 500,
): string {
  if (typeof widthOrSize === 'string') {
    return tmdbImage(movie.poster_path, widthOrSize) ?? FALLBACK_POSTER
  }
  return optimizeCatalogImageUrl(movie.poster_path, {
    width: widthOrSize,
    role: 'poster',
  })
}

export function backdropUrl(
  movie: Pick<CatalogMovie, 'backdrop_path'>,
  widthOrSize: number | TmdbImageSize = 1280,
): string {
  if (typeof widthOrSize === 'string') {
    return tmdbImage(movie.backdrop_path, widthOrSize) ?? FALLBACK_BACKDROP
  }
  return optimizeCatalogImageUrl(movie.backdrop_path, {
    width: widthOrSize,
    role: 'backdrop',
  })
}

export function fallbackPosterDataUrl(): string {
  return FALLBACK_POSTER
}

export function fallbackBackdropDataUrl(): string {
  return FALLBACK_BACKDROP
}
