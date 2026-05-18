import { create } from 'zustand'
import type { TicketCategoryId } from '../constants/ticketCategories'
import type { Hall, Movie, Seat, Session } from '../types'

export type BookingStep = 1 | 2 | 3 | 4

interface BookingState {
  movie: Movie | null
  session: Session | null
  hall: Hall | null
  selectedSeats: Seat[]
  ticketCategory: TicketCategoryId
  step: BookingStep
  setMovie: (m: Movie) => void
  setSession: (s: Session) => void
  setHall: (h: Hall) => void
  toggleSeat: (seat: Seat) => void
  setStep: (n: BookingStep) => void
  setTicketCategory: (c: TicketCategoryId) => void
  clearSeats: () => void
  clear: () => void
}

const initial = {
  movie: null as Movie | null,
  session: null as Session | null,
  hall: null as Hall | null,
  selectedSeats: [] as Seat[],
  ticketCategory: 'adult' as TicketCategoryId,
  step: 3 as BookingStep,
}

export const useBookingStore = create<BookingState>((set, get) => ({
  ...initial,
  setMovie: (movie) => set({ movie }),
  setSession: (session) => set({ session }),
  setHall: (hall) => set({ hall }),
  toggleSeat: (seat) => {
    const { selectedSeats } = get()
    if (!seat.is_available) return
    const exists = selectedSeats.some((s) => s.id === seat.id)
    if (exists) {
      set({ selectedSeats: selectedSeats.filter((s) => s.id !== seat.id) })
      return
    }
    if (selectedSeats.length >= 6) return
    set({ selectedSeats: [...selectedSeats, seat] })
  },
  setStep: (step) => set({ step }),
  setTicketCategory: (ticketCategory) => set({ ticketCategory }),
  clearSeats: () => set({ selectedSeats: [] }),
  clear: () => set({ ...initial }),
}))
