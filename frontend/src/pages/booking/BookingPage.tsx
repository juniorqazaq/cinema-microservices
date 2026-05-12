import { useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { useSessionSeats } from '../../hooks/useSessions'
import { SeatMap } from '../../components/booking/SeatMap'
import { OrderSummary } from '../../components/booking/OrderSummary'
import { PaymentPanel } from '../../components/booking/PaymentPanel'
import { Button } from '../../components/ui/Button'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { useBookingStore } from '../../store/bookingStore'
import { formatDate, formatTime } from '../../utils/format'
import { getErrorMessage } from '../../utils/errorHandler'

const STEP_LABELS = ['Cinema', 'Date & time', 'Seats', 'Payment'] as const

export function BookingPage() {
  const navigate = useNavigate()
  const movie = useBookingStore((s) => s.movie)
  const session = useBookingStore((s) => s.session)
  const hall = useBookingStore((s) => s.hall)
  const step = useBookingStore((s) => s.step)
  const setStep = useBookingStore((s) => s.setStep)
  const selectedSeats = useBookingStore((s) => s.selectedSeats)

  const seatsQuery = useSessionSeats(session?.id)
  const [payBanner, setPayBanner] = useState<string | null>(null)

  if (!movie || !session || !hall) {
    return <Navigate to="/movies" replace />
  }

  const seats = seatsQuery.data ?? []

  function stepClass(n: 1 | 2 | 3 | 4): string {
    const done = (n <= 2 && step >= 3) || (n > 2 && step > n)
    const active = step === n
    if (done) return 'border-accent bg-accentDim text-accent'
    if (active) return 'border-accent bg-accent text-white'
    return 'border-border bg-card2 text-muted'
  }

  return (
    <div className="mx-auto flex max-w-5xl flex-col gap-6">
      <div className="flex flex-wrap gap-2">
        {STEP_LABELS.map((label, idx) => {
          const n = (idx + 1) as 1 | 2 | 3 | 4
          return (
            <button
              key={label}
              type="button"
              onClick={() => {
                if (n === 2) {
                  navigate(`/movies/${movie.id}`)
                  return
                }
                if (n === 3 && step === 4) setStep(3)
              }}
              className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${stepClass(n)}`}
            >
              Step {n}: {label}
            </button>
          )
        })}
      </div>

      {step === 3 ? (
        <>
          {seatsQuery.isLoading ? (
            <div className="flex flex-col gap-2" aria-hidden>
              {Array.from({ length: 6 }).map((_, r) => (
                <div key={r} className="flex gap-1">
                  {Array.from({ length: 10 }).map((__, c) => (
                    <div
                      key={c}
                      className="h-[22px] w-[26px] animate-pulse rounded bg-card2"
                    />
                  ))}
                </div>
              ))}
            </div>
          ) : seatsQuery.isError ? (
            <ErrorBanner message={getErrorMessage(seatsQuery.error)} />
          ) : (
            <div className="flex flex-col gap-6 lg:flex-row lg:items-start">
              <div className="min-w-0 flex-1">
                <SeatMap seats={seats} />
              </div>
              <OrderSummary movie={movie} session={session} hall={hall} />
            </div>
          )}
          <div className="flex flex-wrap justify-end gap-2">
            <Link
              to={`/movies/${movie.id}`}
              className="inline-flex items-center rounded-lg border border-border bg-card2 px-4 py-2 text-body font-medium text-white hover:border-border2"
            >
              Change session
            </Link>
            <Button
              type="button"
              disabled={selectedSeats.length === 0}
              onClick={() => setStep(4)}
            >
              Buy ticket
            </Button>
          </div>
        </>
      ) : null}

      {step === 4 ? (
        <div className="flex flex-col gap-4">
          {payBanner ? (
            <ErrorBanner message={payBanner} onDismiss={() => setPayBanner(null)} />
          ) : null}
          <PaymentPanel
            movie={movie}
            session={session}
            hall={hall}
            selectedSeats={selectedSeats}
            onPaid={({ bookingId }) => {
              const seatLabels = selectedSeats
                .map((s) => `${s.row}${s.number}`)
                .join(', ')
              const totalPaid = selectedSeats.length * session.price
              navigate('/booking/success', {
                replace: true,
                state: {
                  bookingId,
                  movieTitle: movie.title,
                  hallName: hall.name,
                  dateLabel: formatDate(session.start_time),
                  timeLabel: formatTime(session.start_time),
                  seatLabels,
                  totalPaid,
                },
              })
            }}
            onError={(msg) => setPayBanner(msg)}
          />
          <Button type="button" variant="secondary" onClick={() => setStep(3)}>
            Back to seats
          </Button>
        </div>
      ) : null}
    </div>
  )
}
