import { addDays, eachDayOfInterval, format, isSameDay, parseISO } from 'date-fns'
import { useMemo, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useQueries, useQuery } from '@tanstack/react-query'
import axios from 'axios'
import { getHall, getSessions } from '../../api/movies'
import { useBackendMovieId } from '../../hooks/useBackendMovieId'
import { useBookingStore } from '../../store/bookingStore'
import { useAuthStore } from '../../store/authStore'
import { SessionPicker } from './SessionPicker'
import { Spinner } from '../ui/Spinner'
import { mapApiError } from '../../utils/errorHandler'
import type { Hall, Movie, Session } from '../../types'

type PickStep = 'cinema' | 'datetime'

interface ShowtimeBookingPanelProps {
  movie: Movie
}

export function ShowtimeBookingPanel({ movie }: ShowtimeBookingPanelProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const setMovie = useBookingStore((s) => s.setMovie)
  const setSession = useBookingStore((s) => s.setSession)
  const setHall = useBookingStore((s) => s.setHall)
  const setStep = useBookingStore((s) => s.setStep)
  const clearSeats = useBookingStore((s) => s.clearSeats)

  const { backendMovieId, isLoading: resolving, isSynced, backendMovies } =
    useBackendMovieId(movie)

  const [pickStep, setPickStep] = useState<PickStep>('cinema')
  const [selectedHallId, setSelectedHallId] = useState<string | null>(null)
  const [tabIndex, setTabIndex] = useState(0)

  const days = useMemo(() => {
    const start = new Date()
    start.setHours(0, 0, 0, 0)
    return eachDayOfInterval({ start, end: addDays(start, 6) })
  }, [])

  const selectedDate = days[tabIndex] ?? days[0]

  const allSessionsQuery = useQuery({
    queryKey: ['sessions', 'all', backendMovieId],
    queryFn: () => getSessions({ movie_id: backendMovieId }),
    enabled: Boolean(backendMovieId),
  })

  const allSessions = allSessionsQuery.data ?? []

  const hallIds = useMemo(
    () => Array.from(new Set(allSessions.map((s) => s.hall_id))),
    [allSessions],
  )

  const hallQueries = useQueries({
    queries: hallIds.map((hallId) => ({
      queryKey: ['hall', hallId],
      queryFn: () => getHall(hallId),
      enabled: hallIds.length > 0,
    })),
  })

  const hallsById = useMemo(() => {
    const map: Record<string, Hall> = {}
    hallIds.forEach((id, i) => {
      const h = hallQueries[i]?.data
      if (h) map[id] = h
    })
    return map
  }, [hallIds, hallQueries])

  const hallsLoading = hallQueries.some((q) => q.isLoading)

  const sessionsForHallAndDay = useMemo(() => {
    if (!selectedHallId) return []
    return allSessions.filter(
      (s) =>
        s.hall_id === selectedHallId &&
        isSameDay(parseISO(s.start_time), selectedDate),
    )
  }, [allSessions, selectedHallId, selectedDate])

  async function handleSessionPick(session: Session) {
    if (!isAuthenticated) {
      navigate('/login', { state: { from: location } })
      return
    }
    setMovie(movie)
    setSession(session)
    clearSeats()
    const hall = hallsById[session.hall_id] ?? (await getHall(session.hall_id))
    setHall(hall)
    setStep(3)
    navigate('/booking')
  }

  if (resolving) {
    return (
      <section id="showtime-booking" className="scroll-mt-28 rounded-lg border border-border bg-card p-6">
        <Spinner className="h-6 w-6" />
      </section>
    )
  }

  if (!isSynced) {
    return (
      <section
        id="showtime-booking"
        className="scroll-mt-28 rounded-lg border border-border bg-card2/50 p-4"
      >
        <h2 className="text-section font-light text-white">Book tickets</h2>
        <p className="mt-2 text-body text-muted">
          «{movie.title}» is not in the cinema schedule yet. Add it in admin or pick a
          title that already has showtimes.
        </p>
        {backendMovies.length > 0 ? (
          <ul className="mt-4 flex flex-col gap-2">
            {backendMovies.slice(0, 6).map((m) => (
              <li key={m.id}>
                <Link
                  to={`/movies/${m.id}`}
                  className="text-body text-accent hover:underline"
                >
                  {m.title}
                </Link>
              </li>
            ))}
          </ul>
        ) : (
          <p className="mt-3 text-body text-muted">
            No movies in the venue API yet. Create a movie, hall, and sessions in admin.
          </p>
        )}
      </section>
    )
  }

  if (allSessionsQuery.isLoading || hallsLoading) {
    return (
      <section id="showtime-booking" className="scroll-mt-28 rounded-lg border border-border bg-card p-6">
        <h2 className="text-section font-light text-white">Book tickets</h2>
        <div className="mt-4">
          <Spinner className="h-6 w-6" />
        </div>
      </section>
    )
  }

  if (allSessionsQuery.isError) {
    return (
      <section id="showtime-booking" className="scroll-mt-28 rounded-lg border border-border bg-card p-4">
        <p className="text-body text-danger">
          {mapApiError(
            axios.isAxiosError(allSessionsQuery.error)
              ? allSessionsQuery.error.response?.status
              : undefined,
          ) ?? 'Could not load showtimes.'}
        </p>
      </section>
    )
  }

  if (hallIds.length === 0) {
    return (
      <section id="showtime-booking" className="scroll-mt-28 rounded-lg border border-border bg-card2/50 p-4">
        <h2 className="text-section font-light text-white">Book tickets</h2>
        <p className="mt-2 text-body text-muted">
          No showtimes for this title yet. Add sessions in admin.
        </p>
      </section>
    )
  }

  return (
    <section
      id="showtime-booking"
      className="scroll-mt-28 rounded-lg border border-border bg-card p-4 md:p-6"
    >
      <h2 className="text-section font-light text-white">Book tickets</h2>
      <ol className="mt-4 flex flex-wrap gap-2 text-body">
        <StepChip
          n={1}
          label="Cinema"
          active={pickStep === 'cinema'}
          done={Boolean(selectedHallId)}
        />
        <StepChip
          n={2}
          label="Date & time"
          active={pickStep === 'datetime'}
          done={false}
        />
      </ol>

      {pickStep === 'cinema' ? (
        <div className="mt-4">
          <p className="text-body text-muted">Choose a cinema</p>
          <div className="mt-3 grid gap-2 sm:grid-cols-2">
            {hallIds.map((hid) => {
              const hall = hallsById[hid]
              const count = allSessions.filter((s) => s.hall_id === hid).length
              return (
                <button
                  key={hid}
                  type="button"
                  onClick={() => {
                    setSelectedHallId(hid)
                    setPickStep('datetime')
                    setTabIndex(0)
                  }}
                  className="flex flex-col items-start gap-1 rounded-lg border border-border bg-card2 px-4 py-3 text-left transition-colors hover:border-accent"
                >
                  <span className="text-body font-medium text-white">
                    {hall?.name ?? 'Cinema'}
                  </span>
                  <span className="text-[11px] text-muted">
                    {count} showtime{count === 1 ? '' : 's'} this week
                  </span>
                </button>
              )
            })}
          </div>
        </div>
      ) : null}

      {pickStep === 'datetime' && selectedHallId ? (
        <div className="mt-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p className="text-body text-white">
              {hallsById[selectedHallId]?.name ?? 'Cinema'}
            </p>
            <button
              type="button"
              onClick={() => {
                setPickStep('cinema')
                setSelectedHallId(null)
              }}
              className="text-body text-accent hover:underline"
            >
              Change cinema
            </button>
          </div>
          <div className="mt-4 flex flex-wrap gap-2">
            {days.map((d, i) => {
              const active = i === tabIndex
              return (
                <button
                  key={d.toISOString()}
                  type="button"
                  onClick={() => setTabIndex(i)}
                  className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${
                    active
                      ? 'border-accent bg-accentDim text-white'
                      : 'border-border bg-card2 text-muted hover:border-border2'
                  }`}
                >
                  {format(d, 'EEE d')}
                </button>
              )
            })}
          </div>
          <div className="mt-4">
            <SessionPicker
              sessions={sessionsForHallAndDay}
              hallLabel={() => hallsById[selectedHallId]?.name ?? ''}
              onSelect={(s) => void handleSessionPick(s)}
            />
          </div>
          {!isAuthenticated ? (
            <p className="mt-3 text-body text-muted">
              <Link to="/login" state={{ from: location }} className="text-accent hover:underline">
                Sign in
              </Link>{' '}
              to continue to seat selection.
            </p>
          ) : null}
        </div>
      ) : null}
    </section>
  )
}

function StepChip({
  n,
  label,
  active,
  done,
}: {
  n: number
  label: string
  active: boolean
  done: boolean
}) {
  const cls = done
    ? 'border-accent bg-accentDim text-accent'
    : active
      ? 'border-accent bg-accent text-white'
      : 'border-border bg-card2 text-muted'
  return (
    <li
      className={`rounded-lg border px-3 py-1.5 font-medium ${cls}`}
    >
      Step {n}: {label}
    </li>
  )
}

export function scrollToShowtimeBooking() {
  document.getElementById('showtime-booking')?.scrollIntoView({
    behavior: 'smooth',
    block: 'start',
  })
}
