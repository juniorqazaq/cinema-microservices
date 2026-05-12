import { useQuery } from '@tanstack/react-query'
import { getMovie } from '../api/movies'

export function useMovie(
  id: string | undefined,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['movie', id],
    queryFn: () => getMovie(id!),
    enabled: Boolean(id) && (options?.enabled ?? true),
  })
}
