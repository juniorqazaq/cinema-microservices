import type { CatalogMovie, SortCatalogBy } from '../types/movie'
import { filterCatalogByQuery, sortCatalogMovies } from './movieFilters'
import trendingData from '../data/trending.json'
import upcomingData from '../data/upcoming.json'
import nowPlayingData from '../data/nowPlaying.json'
import topRatedData from '../data/topRated.json'

const trending = trendingData as CatalogMovie[]
const upcoming = upcomingData as CatalogMovie[]
const nowPlaying = nowPlayingData as CatalogMovie[]
const topRated = topRatedData as CatalogMovie[]

/** Single deduped list (later files override earlier on id collision). */
export function getAllCatalogMovies(): CatalogMovie[] {
  const map = new Map<number, CatalogMovie>()
  for (const m of trending) map.set(m.id, m)
  for (const m of upcoming) map.set(m.id, m)
  for (const m of nowPlaying) map.set(m.id, m)
  for (const m of topRated) map.set(m.id, m)
  return [...map.values()]
}

export function getTrendingMovies(sortBy: SortCatalogBy = 'popularity'): CatalogMovie[] {
  return sortCatalogMovies([...trending], sortBy, 'desc')
}

export function getUpcomingMovies(sortBy: SortCatalogBy = 'popularity'): CatalogMovie[] {
  return sortCatalogMovies([...upcoming], sortBy, 'desc')
}

export function getNowPlayingMovies(sortBy: SortCatalogBy = 'popularity'): CatalogMovie[] {
  return sortCatalogMovies([...nowPlaying], sortBy, 'desc')
}

export function getTopRatedMovies(sortBy: SortCatalogBy = 'vote_average'): CatalogMovie[] {
  return sortCatalogMovies([...topRated], sortBy, 'desc')
}

export function getMovieById(id: string | number): CatalogMovie | undefined {
  const n = typeof id === 'string' ? Number.parseInt(id, 10) : id
  if (!Number.isFinite(n)) return undefined
  return getAllCatalogMovies().find((m) => m.id === n)
}

export function searchMovies(
  query: string,
  opts?: { sortBy?: SortCatalogBy; order?: 'asc' | 'desc' },
): CatalogMovie[] {
  const base = getAllCatalogMovies()
  const filtered = filterCatalogByQuery(base, query)
  const sortBy = opts?.sortBy ?? 'popularity'
  return sortCatalogMovies(filtered, sortBy, opts?.order ?? 'desc')
}

export function isCatalogMovieId(id: string | undefined): boolean {
  if (!id) return false
  return getMovieById(id) !== undefined
}
