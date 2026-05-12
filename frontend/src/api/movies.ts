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

export async function searchMovies(q: string): Promise<Movie[]> {
  return get<Movie[]>('/movies/search', { params: { q } })
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

/** @deprecated use named exports matching the API spec */
export const fetchMovies = getMovies
export const fetchMovie = getMovie
export const fetchSessions = getSessions
export const fetchSession = getSession
export const fetchSessionSeats = getAvailableSeats
export const fetchHall = getHall
