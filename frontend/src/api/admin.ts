import { del, get, post, put } from './axios'
import type {
  AdminBookingsStats,
  Booking,
  Hall,
  Movie,
  Session,
  User,
} from '../types'

export async function getAllUsers(
  page?: number,
  limit?: number,
): Promise<User[]> {
  return get<User[]>('/admin/users', { params: { page, limit } })
}

export async function banUser(id: string): Promise<Record<string, never>> {
  return post<Record<string, never>>(`/admin/users/${id}/ban`, {})
}

export async function updateUserRole(
  id: string,
  role: 'user' | 'admin',
): Promise<User> {
  return put<User>(`/admin/users/${id}/role`, { role })
}

export async function createMovie(
  data: Omit<Movie, 'id' | 'created_at'>,
): Promise<Movie> {
  return post<Movie>('/admin/movies', data)
}

export async function updateMovie(
  id: string,
  data: Partial<Omit<Movie, 'id' | 'created_at'>>,
): Promise<Movie> {
  return put<Movie>(`/admin/movies/${id}`, data)
}

export async function deleteMovie(id: string): Promise<Record<string, never>> {
  return del<Record<string, never>>(`/admin/movies/${id}`)
}

export async function createHall(
  name: string,
  capacity: number,
  city: string,
  cinema_name: string,
): Promise<Hall> {
  return post<Hall>('/admin/halls', { name, capacity, city, cinema_name })
}

export async function createSession(
  movie_id: string,
  hall_id: string,
  start_time: string,
  price: number,
): Promise<Session> {
  return post<Session>('/admin/sessions', {
    movie_id,
    hall_id,
    start_time,
    price,
  })
}

export async function adminListBookings(): Promise<Booking[]> {
  return get<Booking[]>('/admin/bookings')
}

export async function getBookingStats(): Promise<AdminBookingsStats> {
  return get<AdminBookingsStats>('/admin/bookings/stats')
}
