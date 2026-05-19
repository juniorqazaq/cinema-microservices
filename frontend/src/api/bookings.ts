import { del, get, post } from './axios'
import type { Booking, Payment } from '../types'

export async function createBooking(
  session_id: string,
  seat_id: string,
  ticket_category: string,
  amount: number,
): Promise<Booking> {
  return post<Booking>('/bookings', {
    session_id,
    seat_id,
    ticket_category,
    amount,
  })
}

export async function getBooking(id: string): Promise<Booking> {
  return get<Booking>(`/bookings/${id}`)
}

export async function listUserBookings(): Promise<Booking[]> {
  return get<Booking[]>('/bookings')
}

export async function cancelBooking(id: string): Promise<Record<string, never>> {
  return del<Record<string, never>>(`/bookings/${id}`)
}

export async function getBookingHistory(): Promise<Booking[]> {
  return get<Booking[]>('/bookings/history')
}

export async function confirmPayment(
  booking_id: string,
  amount: number,
  use_wallet = false,
): Promise<Payment> {
  return post<Payment>('/payments/confirm', {
    booking_id,
    amount,
    use_wallet,
  })
}

export async function getPayment(id: string): Promise<Payment> {
  return get<Payment>(`/payments/${id}`)
}
