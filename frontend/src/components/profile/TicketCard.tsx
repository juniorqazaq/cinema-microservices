import { useQuery } from '@tanstack/react-query'
import { parseISO, isFuture } from 'date-fns'
import type { Booking, Session } from '../../types'
import { Badge } from '../ui/Badge'
import { categoryLabel } from '../../constants/ticketCategories'
import { formatDateLocal, formatPrice, formatTimeLocal } from '../../utils/format'
import { getAvailableSeats, getHall, getMovie } from '../../api/movies'
import { Clapperboard } from 'lucide-react'

interface TicketCardProps {
  booking: Booking
  session?: Session
}

function statusBadge(status: Booking['status']) {
  if (status === 'confirmed') {
    return (
      <span className="rounded-md border border-[#166534] bg-[#052210] px-2 py-0.5 text-[11px] font-medium text-[#22c55e]">
        {status}
      </span>
    )
  }
  if (status === 'pending') {
    return (
      <span className="rounded-md border border-accent/40 bg-accentDim px-2 py-0.5 text-[11px] font-medium text-accent">
        {status}
      </span>
    )
  }
  return (
    <span className="rounded-md border border-[#222] bg-[#151515] px-2 py-0.5 text-[11px] font-medium text-muted">
      {status}
    </span>
  )
}

export function TicketCard({ booking, session }: TicketCardProps) {
  const movieQuery = useQuery({
    queryKey: ['movie', session?.movie_id],
    queryFn: () => getMovie(session!.movie_id),
    enabled: Boolean(session?.movie_id),
  })
  const hallQuery = useQuery({
    queryKey: ['hall', session?.hall_id],
    queryFn: () => getHall(session!.hall_id),
    enabled: Boolean(session?.hall_id),
  })
  const seatsQuery = useQuery({
    queryKey: ['session-seats', session?.id],
    queryFn: () => getAvailableSeats(session!.id),
    enabled: Boolean(session?.id),
  })

  const seat = seatsQuery.data?.find((s) => s.id === booking.seat_id)
  const seatLabel = seat ? `${seat.row}${seat.number}` : booking.seat_id.slice(0, 8)

  const start = session?.start_time
  const sessionFuture =
    start && isFuture(parseISO(start)) ? true : start ? false : null

  const shortId =
    booking.id.length > 10 ? `${booking.id.slice(0, 8)}…` : booking.id

  const cancelled = booking.status === 'cancelled'
  const movieTitle = movieQuery.data?.title ?? 'Movie'
  const cinemaName =
    hallQuery.data?.cinema_name ?? hallQuery.data?.name ?? session?.cinema_name ?? 'Cinema'
  const qrUrl = `https://api.qrserver.com/v1/create-qr-code/?size=120x120&data=${encodeURIComponent(booking.id)}`

  return (
    <div
      className={`flex gap-4 rounded-lg border border-border bg-card p-4 ${
        cancelled ? 'opacity-60' : ''
      }`}
    >
      <div className="flex shrink-0 flex-col items-center gap-2">
        <div className="flex h-14 w-10 items-center justify-center rounded border border-border bg-card2">
          <Clapperboard className="h-5 w-5 text-muted" aria-hidden />
        </div>
        {booking.status === 'confirmed' ? (
          <img
            src={qrUrl}
            alt={`Ticket QR ${shortId}`}
            className="h-[72px] w-[72px] rounded border border-border bg-white p-1"
            width={72}
            height={72}
          />
        ) : null}
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <span className="font-mono text-[11px] text-muted">{shortId}</span>
          {statusBadge(booking.status)}
          {booking.status === 'confirmed' && sessionFuture === true ? (
            <Badge tone="warning">Upcoming</Badge>
          ) : null}
          {booking.status === 'confirmed' && sessionFuture === false ? (
            <Badge tone="neutral">Watched</Badge>
          ) : null}
        </div>
        <p className="mt-2 text-body font-medium text-white">{movieTitle}</p>
        <p className="mt-1 text-body text-muted">{cinemaName}</p>
        <p className="mt-2 text-body text-white">
          {start
            ? `${formatDateLocal(start)} · ${formatTimeLocal(start)}`
            : '—'}
        </p>
        <p className="mt-1 text-body text-muted">Seat {seatLabel}</p>
        {booking.ticket_category ? (
          <p className="mt-1 text-body text-muted">
            {categoryLabel(booking.ticket_category)}
            {booking.amount_paid != null && booking.amount_paid > 0
              ? ` · ${formatPrice(booking.amount_paid)}`
              : ''}
          </p>
        ) : null}
      </div>
    </div>
  )
}
