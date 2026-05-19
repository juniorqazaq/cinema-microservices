import type { Hall, Movie, Session } from '../../types'
import { formatDate, formatPrice, formatTimeLocal } from '../../utils/format'
import { useBookingStore } from '../../store/bookingStore'
import { ticketPrice } from '../../constants/ticketCategories'

interface OrderSummaryProps {
  movie: Movie
  session: Session
  hall: Hall
}

function isVipRow(row: string): boolean {
  const r = row.toUpperCase()
  return r === 'A' || r === 'B'
}

export function OrderSummary({ movie, session, hall }: OrderSummaryProps) {
  const selectedSeats = useBookingStore((s) => s.selectedSeats)
  const ticketCategory = useBookingStore((s) => s.ticketCategory)

  const unitPrice = ticketPrice(session.price, ticketCategory)
  const total = selectedSeats.length * unitPrice

  return (
    <aside className="w-full shrink-0 border border-border bg-card p-3 lg:w-[180px]">
      <p className="text-card-title font-medium text-white">{movie.title}</p>
      <p className="mt-2 text-body text-muted">{hall.name}</p>
      <p className="mt-1 text-body text-white">
        {formatTimeLocal(session.start_time)} · {formatDate(session.start_time)}
      </p>
      <div className="my-3 border-t border-border" />
      {selectedSeats.length === 0 ? (
        <p className="text-body text-muted">No seats selected.</p>
      ) : (
        <ul className="flex flex-col gap-2">
          {selectedSeats.map((seat) => (
            <li key={seat.id} className="text-body text-muted">
              {seat.row}
              {seat.number} · {isVipRow(seat.row) ? 'VIP' : 'Standard'} ·{' '}
              {formatPrice(unitPrice)}
            </li>
          ))}
        </ul>
      )}
      <div className="my-3 border-t border-border" />
      <p className="text-body font-medium text-white">
        Total: {formatPrice(total)}
      </p>
    </aside>
  )
}
