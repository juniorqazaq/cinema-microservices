import { useQuery } from '@tanstack/react-query'
import { format, parseISO } from 'date-fns'
import { isTmdbConfigured, searchTmdbMovie } from '../api/tmdb'
import type { Movie } from '../types'

/** Build query: "Title" or "Title 2023" from backend movie. */
export function tmdbSearchQuery(movie: Pick<Movie, 'title' | 'created_at'>): string {
  let q = movie.title.trim()
  if (movie.created_at) {
    try {
      q = `${q} ${format(parseISO(movie.created_at), 'yyyy')}`
    } catch {
      /* ignore */
    }
  }
  return q
}

export function useTmdbMovie(movie: Pick<Movie, 'title' | 'created_at'> | undefined) {
  const q = movie ? tmdbSearchQuery(movie) : ''
  const enabled = Boolean(movie && isTmdbConfigured() && q)

  return useQuery({
    queryKey: ['tmdb', 'search', q],
    queryFn: () => searchTmdbMovie(q),
    enabled,
    staleTime: 1000 * 60 * 60 * 24,
    gcTime: 1000 * 60 * 60 * 48,
    retry: 1,
  })
}
