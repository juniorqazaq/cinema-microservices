import { useQuery } from '@tanstack/react-query'
import {
  getAvailableSeats,
  getSession,
  getSessions,
} from '../api/movies'
import type { SessionsListParams } from '../types'

export function useSessions(
  params?: SessionsListParams,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['sessions', params],
    queryFn: () => getSessions(params!),
    enabled:
      (options?.enabled ?? true) &&
      Boolean(params?.movie_id) &&
      (params?.date !== undefined ? Boolean(params.date) : true),
  })
}

export function useSession(id: string | undefined) {
  return useQuery({
    queryKey: ['session', id],
    queryFn: () => getSession(id!),
    enabled: Boolean(id),
  })
}

export function useSessionSeats(id: string | undefined) {
  return useQuery({
    queryKey: ['session', id, 'seats'],
    queryFn: () => getAvailableSeats(id!),
    enabled: Boolean(id),
  })
}
