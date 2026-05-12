import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '../store/authStore'
import {
  adminListBookings,
  getBookingStats,
  getAllUsers,
} from '../api/admin'
import {
  confirmPayment,
  createBooking,
  getBookingHistory,
} from '../api/bookings'

export function useBookingHistory() {
  const enabled = useAuthStore((s) => s.isAuthenticated)
  return useQuery({
    queryKey: ['bookings', 'history'],
    queryFn: getBookingHistory,
    enabled,
  })
}

export function useAdminStats() {
  return useQuery({
    queryKey: ['admin', 'bookings', 'stats'],
    queryFn: getBookingStats,
  })
}

export function useAdminBookings() {
  return useQuery({
    queryKey: ['admin', 'bookings'],
    queryFn: adminListBookings,
  })
}

export function useAdminUsers(page = 1, limit = 50) {
  return useQuery({
    queryKey: ['admin', 'users', page, limit],
    queryFn: () => getAllUsers(page, limit),
  })
}

export function useBookSeatAndPay() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (vars: {
      session_id: string
      seat_id: string
      amount: number
    }) => {
      const booking = await createBooking(vars.session_id, vars.seat_id)
      const payment = await confirmPayment(booking.id, vars.amount)
      return { booking, payment }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bookings'] })
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })
}
