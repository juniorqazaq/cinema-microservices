import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getMovies } from '../api/movies'
import type { Movie } from '../types'

const UUID_RE =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

function normTitle(s: string): string {
  return s.trim().toLowerCase().replace(/\s+/g, ' ')
}

/** Map catalog / display movie to backend UUID when synced. */
export function useBackendMovieId(displayMovie: Movie | undefined) {
  const isBackendId = Boolean(displayMovie?.id && UUID_RE.test(displayMovie.id))

  const moviesQuery = useQuery({
    queryKey: ['movies', 'all'],
    queryFn: () => getMovies({ limit: 200 }),
    enabled: Boolean(displayMovie) && !isBackendId,
    staleTime: 60_000,
  })

  const backendMovieId = useMemo(() => {
    if (!displayMovie) return undefined
    if (isBackendId) return displayMovie.id
    const list = moviesQuery.data ?? []
    const title = normTitle(displayMovie.title)
    const match = list.find((m) => normTitle(m.title) === title)
    return match?.id
  }, [displayMovie, isBackendId, moviesQuery.data])

  return {
    backendMovieId,
    isBackendId,
    isLoading: !isBackendId && moviesQuery.isLoading,
    isSynced: Boolean(backendMovieId),
    backendMovies: moviesQuery.data ?? [],
  }
}
