import { useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { CircleCheck } from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { formatPrice } from '../../utils/format'

interface SuccessState {
  bookingId?: string
  movieTitle?: string
  hallName?: string
  dateLabel?: string
  timeLabel?: string
  seatLabels?: string
  totalPaid?: number
}

function bookingRef(id: string): string {
  const alnum = id.replace(/[^a-zA-Z0-9]/g, '').toUpperCase()
  const tail = alnum.slice(-4).padStart(4, '0')
  return `#CM-${tail}`
}

export function BookingSuccessPage() {
  const location = useLocation()
  const clear = useBookingStore((s) => s.clear)
  const state = (location.state ?? {}) as SuccessState

  useEffect(() => {
    clear()
  }, [clear])

  const title = state.movieTitle ?? 'Your booking'
  const hall = state.hallName ?? '—'
  const date = state.dateLabel ?? '—'
  const time = state.timeLabel ?? '—'
  const seats = state.seatLabels ?? '—'
  const total = state.totalPaid ?? 0
  const id = state.bookingId ?? ''

  return (
    <div className="mx-auto flex max-w-lg flex-col items-center gap-6 text-center">
      <CircleCheck className="h-[48px] w-[48px] text-[#22c55e]" aria-hidden />
      <div>
        <h2 className="text-section font-light text-white">Booking confirmed!</h2>
        <p className="mt-2 text-body text-muted">
          Your tickets have been sent to your email.
        </p>
      </div>
      <div className="w-full rounded-lg border border-border bg-card p-4 text-left">
        <p className="text-card-title font-medium text-white">{title}</p>
        <p className="mt-2 text-body text-muted">
          {hall} · {date} · {time}
        </p>
        <p className="mt-2 text-body text-white">Seats: {seats}</p>
        <p className="mt-2 text-body text-muted">
          Total paid: {formatPrice(total)}
        </p>
        <div className="my-4 border-t border-dashed border-border" />
        <p className="text-body text-muted">
          Booking ID: {id ? bookingRef(id) : '—'}
        </p>
      </div>
      <Link
        to="/profile"
        className="inline-flex w-full max-w-xs items-center justify-center gap-2 rounded-lg border border-accent bg-accent px-4 py-2 text-body font-medium text-white transition-colors duration-150 hover:border-accentHover hover:bg-accentHover"
      >
        View my tickets
      </Link>
    </div>
  )
}
