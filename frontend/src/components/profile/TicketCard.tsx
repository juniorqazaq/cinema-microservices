import { parseISO, isFuture } from 'date-fns'
import type { Booking, Session } from '../../types'
import { Badge } from '../ui/Badge'
import { formatDate, formatTime } from '../../utils/format'
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
  const start = session?.start_time
  const sessionFuture =
    start && isFuture(parseISO(start)) ? true : start ? false : null

  const shortId =
    booking.id.length > 10 ? `${booking.id.slice(0, 8)}…` : booking.id

  const cancelled = booking.status === 'cancelled'

  return (
    <div
      className={`flex gap-3 rounded-lg border border-border bg-card p-3 ${
        cancelled ? 'opacity-60' : ''
      }`}
    >
      <div className="flex h-14 w-10 shrink-0 items-center justify-center rounded border border-border bg-card2">
        <Clapperboard className="h-5 w-5 text-muted" aria-hidden />
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
        <p className="mt-2 text-body text-white">
          {start
            ? `${formatDate(start)} · ${formatTime(start)}`
            : booking.session_id}
        </p>
        <p className="mt-1 text-body text-muted">Seat: {booking.seat_id}</p>
      </div>
    </div>
  )
}
