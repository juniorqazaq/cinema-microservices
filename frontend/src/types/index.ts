export interface User {
  id: string
  email: string
  role: 'user' | 'admin'
  is_banned: boolean
  created_at: string
  updated_at: string
  /** Optional wallet balance when present in JWT or profile */
  balance?: number
}

export interface Movie {
  id: string
  title: string
  description: string
  genre: string
  duration: number
  rating: number
  created_at: string
}

export interface Hall {
  id: string
  name: string
  capacity: number
}

export interface Seat {
  id: string
  hall_id: string
  row: string
  number: number
  is_available: boolean
}

export interface Session {
  id: string
  movie_id: string
  hall_id: string
  start_time: string
  price: number
}

export interface Booking {
  id: string
  user_id: string
  session_id: string
  seat_id: string
  status: 'pending' | 'confirmed' | 'cancelled'
  created_at: string
}

export interface Payment {
  id: string
  booking_id: string
  amount: number
  status: 'pending' | 'paid' | 'refunded'
  paid_at: string | null
}

export type ApiSuccess<T> = { data: T }
export type ApiErrorBody = { error: string }

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface RefreshResponse {
  access_token: string
  refresh_token: string
}

export interface MoviesListParams {
  page?: number
  limit?: number
  genre?: string
}

export interface SessionsListParams {
  movie_id?: string
  date?: string
}

export interface AdminBookingsStats {
  total: number
  confirmed: number
  cancelled: number
}
