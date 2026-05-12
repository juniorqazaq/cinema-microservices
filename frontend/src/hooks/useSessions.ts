import { useQuery } from '@tanstack/react-query'
import { fetchSession, fetchSessionSeats, fetchSessions } from '../api/movies'
import type { SessionsListParams } from '../types'

export function useSessions(
  params?: SessionsListParams,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['sessions', params],
    queryFn: () => fetchSessions(params!),
    enabled:
      (options?.enabled ?? true) &&
      Boolean(params?.movie_id) &&
      Boolean(params?.date),
  })
}

export function useSession(id: string | undefined) {
  return useQuery({
    queryKey: ['session', id],
    queryFn: () => fetchSession(id!),
    enabled: Boolean(id),
  })
}

export function useSessionSeats(id: string | undefined) {
  return useQuery({
    queryKey: ['session', id, 'seats'],
    queryFn: () => fetchSessionSeats(id!),
    enabled: Boolean(id),
  })
}
