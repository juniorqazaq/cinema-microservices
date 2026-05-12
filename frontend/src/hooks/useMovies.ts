import { useQuery } from '@tanstack/react-query'
import { fetchMovie, fetchMovies, searchMovies } from '../api/movies'
import type { MoviesListParams } from '../types'

export function useMoviesList(
  params?: MoviesListParams,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['movies', params],
    queryFn: () => fetchMovies(params),
    enabled: options?.enabled ?? true,
  })
}

export function useMovie(
  id: string | undefined,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['movie', id],
    queryFn: () => fetchMovie(id!),
    enabled: Boolean(id) && (options?.enabled ?? true),
  })
}

export function useMovieSearch(q: string) {
  return useQuery({
    queryKey: ['movies', 'search', q],
    queryFn: () => searchMovies(q),
    enabled: q.trim().length > 0,
  })
}
