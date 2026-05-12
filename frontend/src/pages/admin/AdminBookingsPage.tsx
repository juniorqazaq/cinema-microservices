import { useMemo, useState } from 'react'
import { format, parseISO } from 'date-fns'
import { useAdminBookings } from '../../hooks/useBooking'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import { getErrorMessage } from '../../utils/errorHandler'
import type { Booking } from '../../types'

function bookingStatusClass(status: Booking['status']): string {
  if (status === 'pending')
    return 'border-accent/40 bg-accentDim text-accent'
  if (status === 'confirmed') return 'border-[#166534] bg-[#052210] text-[#22c55e]'
  return 'border-[#222] bg-[#151515] text-muted'
}

export function AdminBookingsPage() {
  const [status, setStatus] = useState<string>('')
  const bookingsQuery = useAdminBookings()

  const rows = useMemo(() => bookingsQuery.data ?? [], [bookingsQuery.data])

  const filtered = useMemo(() => {
    if (!status) return rows
    return rows.filter((b) => b.status === status)
  }, [rows, status])

  if (bookingsQuery.isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner className="h-6 w-6" />
      </div>
    )
  }

  if (bookingsQuery.isError) {
    return <ErrorBanner message={getErrorMessage(bookingsQuery.error)} />
  }

  const filters = [
    { value: '', label: 'All' },
    { value: 'pending', label: 'pending' },
    { value: 'confirmed', label: 'confirmed' },
    { value: 'cancelled', label: 'cancelled' },
  ]

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <h1 className="text-section font-light text-white">Bookings</h1>
        <div className="flex flex-wrap gap-2">
          {filters.map((f) => (
            <button
              key={f.label}
              type="button"
              onClick={() => setStatus(f.value)}
              className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${
                status === f.value
                  ? 'border-accent bg-accentDim text-accent'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>
      <div className="overflow-x-auto rounded-lg border border-border bg-card">
        <table className="min-w-full divide-y divide-border text-left text-body">
          <thead className="text-[11px] font-medium text-muted">
            <tr>
              <th className="px-4 py-2">Id</th>
              <th className="px-4 py-2">User</th>
              <th className="px-4 py-2">Session</th>
              <th className="px-4 py-2">Seat</th>
              <th className="px-4 py-2">Status</th>
              <th className="px-4 py-2">Created</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border text-white">
            {filtered.map((b) => {
              const short = (id: string) =>
                id.length > 10 ? `${id.slice(0, 6)}…` : id
              return (
                <tr key={b.id}>
                  <td className="px-4 py-2 font-mono text-[11px]">{short(b.id)}</td>
                  <td className="px-4 py-2 font-mono text-[11px]">
                    {short(b.user_id)}
                  </td>
                  <td className="px-4 py-2 font-mono text-[11px]">{b.session_id}</td>
                  <td className="px-4 py-2 font-mono text-[11px]">{b.seat_id}</td>
                  <td className="px-4 py-2">
                    <span
                      className={`inline-flex rounded-md border px-2 py-0.5 text-[11px] font-medium ${bookingStatusClass(b.status)}`}
                    >
                      {b.status}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-muted">
                    {format(parseISO(b.created_at), 'PPp')}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
      {filtered.length === 0 ? (
        <p className="text-body text-muted">No bookings match this filter.</p>
      ) : null}
    </div>
  )
}
