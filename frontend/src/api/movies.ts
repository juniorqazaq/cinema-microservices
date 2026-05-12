import { get } from './axios'
import type {
  Hall,
  Movie,
  Seat,
  Session,
  MoviesListParams,
  SessionsListParams,
} from '../types'

export async function getMovies(
  params?: MoviesListParams,
): Promise<Movie[]> {
  return get<Movie[]>('/movies', { params })
}

export async function getMovie(id: string): Promise<Movie> {
  return get<Movie>(`/movies/${id}`)
}

export async function getSessions(
  params?: SessionsListParams,
): Promise<Session[]> {
  return get<Session[]>('/sessions', { params })
}

export async function getSession(id: string): Promise<Session> {
  return get<Session>(`/sessions/${id}`)
}

export async function getAvailableSeats(sessionId: string): Promise<Seat[]> {
  return get<Seat[]>(`/sessions/${sessionId}/seats`)
}

export async function getHall(id: string): Promise<Hall> {
  return get<Hall>(`/halls/${id}`)
}
