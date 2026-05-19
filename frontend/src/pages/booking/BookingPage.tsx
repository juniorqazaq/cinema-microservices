import { useMemo, useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useSessionSeats } from '../../hooks/useSessions'
import { SeatMap } from '../../components/booking/SeatMap'
import { OrderSummary } from '../../components/booking/OrderSummary'
import { PaymentPanel } from '../../components/booking/PaymentPanel'
import { BookingDateTimeBar } from '../../components/booking/BookingDateTimeBar'
import { Button } from '../../components/ui/Button'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { useBookingStore } from '../../store/bookingStore'
import { formatDate, formatPrice, formatTimeLocal } from '../../utils/format'
import { getErrorMessage } from '../../utils/errorHandler'
import { getHall, getSessions } from '../../api/movies'
import { ticketPrice } from '../../constants/ticketCategories'
import type { Session } from '../../types'

const STEP_LABELS = ['Cinema', 'Date & time', 'Seats', 'Payment'] as const

export function BookingPage() {
  const navigate = useNavigate()
  const movie = useBookingStore((s) => s.movie)
  const session = useBookingStore((s) => s.session)
  const hall = useBookingStore((s) => s.hall)
  const step = useBookingStore((s) => s.step)
  const setStep = useBookingStore((s) => s.setStep)
  const selectedSeats = useBookingStore((s) => s.selectedSeats)
  const ticketCategory = useBookingStore((s) => s.ticketCategory)
  const setSession = useBookingStore((s) => s.setSession)
  const setHall = useBookingStore((s) => s.setHall)
  const clearSeats = useBookingStore((s) => s.clearSeats)

  const [tabIndex, setTabIndex] = useState(0)
  const [payBanner, setPayBanner] = useState<string | null>(null)
  const [switchingSession, setSwitchingSession] = useState(false)

  const seatsQuery = useSessionSeats(session?.id)

  const sessionsQuery = useQuery({
    queryKey: ['sessions', movie?.id, session?.city, session?.cinema_name],
    queryFn: () =>
      getSessions({
        movie_id: movie!.id,
        city: session?.city,
      }),
    enabled: Boolean(movie?.id && session?.city && step === 3),
  })

  const cinemaSessions = useMemo(() => {
    const list = sessionsQuery.data ?? []
    const cinema = session?.cinema_name
    if (!cinema) return list
    return list.filter((s) => (s.cinema_name || s.hall_name) === cinema)
  }, [sessionsQuery.data, session?.cinema_name, session?.hall_name])

  if (!movie || !session || !hall) {
    return <Navigate to="/movies" replace />
  }

  const seats = seatsQuery.data ?? []
  const unitPrice = ticketPrice(session.price, ticketCategory)
  const total = selectedSeats.length * unitPrice

  async function handleSessionChange(s: Session) {
    if (!session || s.id === session.id) return
    setSwitchingSession(true)
    clearSeats()
    setSession(s)
    try {
      const h = await getHall(s.hall_id)
      setHall(h)
    } finally {
      setSwitchingSession(false)
    }
  }

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
                if (n <= 2) {
                  navigate(`/movies/${movie.id}/sessions`)
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
          {sessionsQuery.data && cinemaSessions.length > 0 ? (
            <BookingDateTimeBar
              sessions={cinemaSessions}
              selectedSessionId={session.id}
              tabIndex={tabIndex}
              onTabIndexChange={(i) => {
                setTabIndex(i)
              }}
              onSelectSession={(s) => void handleSessionChange(s)}
            />
          ) : null}

          {seatsQuery.isLoading || switchingSession ? (
            <div className="flex flex-col gap-2" aria-hidden>
              {Array.from({ length: 6 }).map((_, r) => (
                <div key={r} className="flex gap-1">
                  {Array.from({ length: 10 }).map((__, c) => (
                    <div
                      key={c}
                      className="h-5 w-5 animate-pulse rounded bg-[#2A2A3D]"
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

          <div className="flex flex-wrap items-center justify-end gap-3">
            <Link
              to={`/movies/${movie.id}/sessions`}
              className="inline-flex items-center rounded-lg border border-border bg-card2 px-4 py-2 text-body font-medium text-white hover:border-border2"
            >
              Change cinema
            </Link>
            <button
              type="button"
              disabled={selectedSeats.length === 0}
              onClick={() => setStep(4)}
              className="inline-flex items-center gap-2 rounded-xl px-8 py-3 text-base font-semibold text-white transition-opacity disabled:cursor-not-allowed disabled:opacity-40"
              style={{
                background: 'linear-gradient(135deg, #7C4DFF 0%, #5b2fd4 50%, #9d6fff 100%)',
                boxShadow: '0 8px 32px rgba(124, 77, 255, 0.45)',
              }}
            >
              Book Ticket Now {formatPrice(total)}
            </button>
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
            onPaid={({ bookingId, totalPaid }) => {
              const seatLabels = selectedSeats
                .map((s) => `${s.row}${s.number}`)
                .join(', ')
              navigate('/booking/success', {
                replace: true,
                state: {
                  bookingId,
                  movieTitle: movie.title,
                  hallName: hall.name,
                  dateLabel: formatDate(session.start_time),
                  timeLabel: formatTimeLocal(session.start_time),
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
