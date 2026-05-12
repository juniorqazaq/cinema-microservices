import type { CatalogMovie } from '../types/movie'
import type { Movie } from '../types'

const GENRE_BY_ID: Record<number, string> = {
  28: 'Action',
  12: 'Adventure',
  16: 'Animation',
  35: 'Comedy',
  80: 'Crime',
  99: 'Documentary',
  18: 'Drama',
  10751: 'Family',
  14: 'Fantasy',
  36: 'History',
  27: 'Horror',
  10402: 'Music',
  9648: 'Mystery',
  10749: 'Romance',
  878: 'Sci-Fi',
  10770: 'TV Movie',
  53: 'Thriller',
  10752: 'War',
  37: 'Western',
}

export function genreIdsToLabel(ids: number[]): string {
  if (!ids.length) return 'Drama'
  return ids
    .slice(0, 3)
    .map((i) => GENRE_BY_ID[i] ?? 'Drama')
    .join(', ')
}

/** Map catalog row → existing `Movie` shape for grids / booking store compatibility. */
export function catalogToMovie(c: CatalogMovie): Movie {
  const duration = c.runtime > 0 ? c.runtime : 120
  return {
    id: String(c.id),
    title: c.title,
    description: c.overview,
    genre: genreIdsToLabel(c.genre_ids),
    duration,
    rating: c.vote_average,
    created_at: c.release_date
      ? `${c.release_date}T12:00:00.000Z`
      : new Date().toISOString(),
  }
}

export function buildCatalogPosterMap(movies: CatalogMovie[]): Record<string, string> {
  const map: Record<string, string> = {}
  for (const m of movies) {
    map[String(m.id)] = m.poster_path
  }
  return map
}
