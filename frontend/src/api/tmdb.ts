import axios from 'axios'
import { tmdbImage, type TmdbImageSize } from '../lib/movieImages'

const TMDB = 'https://api.themoviedb.org/3'

export interface TmdbSearchMovie {
  id: number
  title: string
  poster_path: string | null
  backdrop_path: string | null
  vote_average: number
  release_date?: string
}

function readToken(): string | undefined {
  const t = import.meta.env.VITE_TMDB_READ_ACCESS_TOKEN
  return typeof t === 'string' && t.trim() ? t.trim() : undefined
}

function apiKey(): string | undefined {
  const k = import.meta.env.VITE_TMDB_API_KEY
  return typeof k === 'string' && k.trim() ? k.trim() : undefined
}

export function isTmdbConfigured(): boolean {
  return Boolean(readToken() || apiKey())
}

/** Poster / backdrop URL or null if path missing (uses shared `tmdbImage` helper). */
export function tmdbImg(
  path: string | null | undefined,
  size: TmdbImageSize = 'w500',
): string | null {
  return tmdbImage(path, size)
}

/** First search hit. TMDB v3: Bearer read token OR api_key (see .env). */
export async function searchTmdbMovie(
  query: string,
): Promise<TmdbSearchMovie | null> {
  if (!query.trim()) return null
  const token = readToken()
  const key = apiKey()
  if (!token && !key) return null

  const params: Record<string, string | boolean> = {
    query: query.trim(),
    include_adult: false,
    language: 'en-US',
  }
  if (key && !token) {
    params.api_key = key
  }

  const { data } = await axios.get<{ results: TmdbSearchMovie[] }>(
    `${TMDB}/search/movie`,
    {
      params,
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      timeout: 12_000,
    },
  )
  return data.results?.[0] ?? null
}
