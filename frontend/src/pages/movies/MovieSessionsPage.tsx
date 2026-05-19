import { addDays, eachDayOfInterval } from 'date-fns'
import { useMemo, useState } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getMovieById } from '../../lib/movieService'
import { catalogToMovie } from '../../lib/catalogAdapter'
import { useMovie } from '../../hooks/useMovies'
import { useBackendMovieId } from '../../hooks/useBackendMovieId'
import { getHall, getSessions } from '../../api/movies'
import { useBookingStore } from '../../store/bookingStore'
import { useAuthStore } from '../../store/authStore'
import { DEFAULT_CITY, CITIES, type CityName } from '../../constants/cities'
import { Badge } from '../../components/ui/Badge'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import {
  formatPrice,
  formatTimeLocal,
  localDayKey,
  sessionLocalDayKey,
} from '../../utils/format'
import type { Movie, Session } from '../../types'
import { AgeVerifyModal } from '../../components/booking/AgeVerifyModal'
import { CinemaCard, type CinemaInfo } from '../../components/booking/CinemaCard'

type FlowStep = 1 | 2

function groupCinemas(sessions: Session[], city: string): CinemaInfo[] {
  const map = new Map<string, Set<string>>()
  for (const s of sessions) {
    const name = s.cinema_name || s.hall_name || 'Cinema'
    const halls = map.get(name) ?? new Set<string>()
    halls.add(s.hall_id)
    map.set(name, halls)
  }
  return Array.from(map.entries()).map(([name, halls]) => ({
    name,
    city,
    hallCount: halls.size,
  }))
}

export function MovieSessionsPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const catalogMovie = useMemo(() => (id ? getMovieById(id) : undefined), [id])
  const displayMovie: Movie | undefined = catalogMovie
    ? catalogToMovie(catalogMovie)
    : undefined
  const movieQuery = useMovie(id, { enabled: Boolean(id) && !catalogMovie })
  const movie = displayMovie ?? movieQuery.data

  const { backendMovieId, isLoading: resolving, isSynced } = useBackendMovieId(
    movie ?? {
      id: '',
      title: '',
      description: '',
      genre: '',
      duration: 0,
      rating: 0,
      created_at: '',
    },
  )

  const [flowStep, setFlowStep] = useState<FlowStep>(1)
  const [city, setCity] = useState<CityName>(DEFAULT_CITY)
  const [selectedCinema, setSelectedCinema] = useState<string | null>(null)
  const [tabIndex, setTabIndex] = useState(0)
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null)
  const [ageOpen, setAgeOpen] = useState(false)
  const [pendingSession, setPendingSession] = useState<Session | null>(null)

  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const setMovie = useBookingStore((s) => s.setMovie)
  const setSession = useBookingStore((s) => s.setSession)
  const setHall = useBookingStore((s) => s.setHall)
  const clearSeats = useBookingStore((s) => s.clearSeats)
  const setStep = useBookingStore((s) => s.setStep)

  const days = useMemo(() => {
    const start = new Date()
    start.setHours(0, 0, 0, 0)
    return eachDayOfInterval({ start, end: addDays(start, 6) })
  }, [])

  const sessionsQuery = useQuery({
    queryKey: ['sessions', backendMovieId, city],
    queryFn: () => getSessions({ movie_id: backendMovieId, city }),
    enabled: Boolean(backendMovieId),
  })

  const cinemas = useMemo(
    () => groupCinemas(sessionsQuery.data ?? [], city),
    [sessionsQuery.data, city],
  )

  const cinemaSessions = useMemo(() => {
    const list = sessionsQuery.data ?? []
    if (!selectedCinema) return list
    return list.filter((s) => (s.cinema_name || s.hall_name || 'Cinema') === selectedCinema)
  }, [sessionsQuery.data, selectedCinema])

  const selectedDayKey = localDayKey(days[tabIndex] ?? days[0]!)

  const sessionsForDay = useMemo(
    () =>
      cinemaSessions.filter(
        (s) => sessionLocalDayKey(s.start_time) === selectedDayKey,
      ),
    [cinemaSessions, selectedDayKey],
  )

  const ageRating = movie?.age_rating ?? cinemaSessions[0]?.age_rating ?? 12
  const is18Plus = ageRating >= 18

  async function proceedToBooking(session: Session) {
    if (!movie) return
    if (!isAuthenticated) {
      navigate('/login', { state: { from: location } })
      return
    }
    setMovie(movie)
    setSession(session)
    clearSeats()
    const hall = await getHall(session.hall_id)
    setHall(hall)
    setStep(3)
    navigate('/booking')
  }

  function onPickSession(session: Session) {
    setSelectedSessionId(session.id)
    if (is18Plus) {
      setPendingSession(session)
      setAgeOpen(true)
      return
    }
    void proceedToBooking(session)
  }

  if (!movie && movieQuery.isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    )
  }

  if (!movie) {
    return (
      <div className="rounded-lg border border-border bg-card p-6 text-muted">
        Film not found
      </div>
    )
  }

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <Link
        to={`/movies/${id}`}
        className="inline-flex text-body text-muted hover:text-accent"
      >
        ← Back to {movie.title}
      </Link>

      <header>
        <h1 className="text-section font-light text-white">Book tickets</h1>
        <p className="mt-1 text-body text-muted">{movie.title}</p>
        <div className="mt-3 flex flex-wrap gap-2">
          <StepChip n={1} label="Cinema" active={flowStep === 1} done={flowStep > 1} />
          <StepChip n={2} label="Date & time" active={flowStep === 2} done={false} />
        </div>
        {is18Plus ? (
          <Badge tone="danger" className="mt-2">
            18+
          </Badge>
        ) : null}
      </header>

      <section className="rounded-lg border border-border bg-card p-4">
        <p className="text-body font-medium text-white">City</p>
        <div className="mt-2 flex flex-wrap gap-2">
          {CITIES.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => {
                setCity(c)
                setSelectedCinema(null)
                setFlowStep(1)
              }}
              className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors ${
                city === c
                  ? 'border-accent bg-accentDim text-accent'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              {c}
            </button>
          ))}
        </div>
      </section>

      {resolving ? (
        <Spinner />
      ) : !isSynced ? (
        <p className="text-body text-muted">
          «{movie.title}» is not in the cinema schedule yet. Ask admin to add the movie and
          sessions for {city}.
        </p>
      ) : sessionsQuery.isLoading ? (
        <Spinner />
      ) : sessionsQuery.isError ? (
        <ErrorBanner message="Could not load showtimes." />
      ) : flowStep === 1 ? (
        <section className="flex flex-col gap-3">
          <p className="text-body font-medium text-white">Choose cinema</p>
          {cinemas.length === 0 ? (
            <p className="text-body text-muted">No cinemas in {city} for this movie.</p>
          ) : (
            <div className="grid gap-3 sm:grid-cols-2">
              {cinemas.map((cinema) => (
                <CinemaCard
                  key={cinema.name}
                  cinema={cinema}
                  selected={selectedCinema === cinema.name}
                  onSelect={() => {
                    setSelectedCinema(cinema.name)
                    setFlowStep(2)
                    setTabIndex(0)
                    setSelectedSessionId(null)
                  }}
                />
              ))}
            </div>
          )}
        </section>
      ) : (
        <section className="flex flex-col gap-4">
          <div className="flex items-center justify-between gap-2">
            <p className="text-body font-medium text-white">{selectedCinema}</p>
            <button
              type="button"
              className="text-body text-accent hover:underline"
              onClick={() => setFlowStep(1)}
            >
              Change cinema
            </button>
          </div>

          <div className="flex gap-2 overflow-x-auto pb-1">
            {days.map((d, i) => {
              const dayNum = new Intl.DateTimeFormat('en', {
                timeZone: 'Asia/Almaty',
                day: 'numeric',
              }).format(d)
              const weekday = new Intl.DateTimeFormat('en', {
                timeZone: 'Asia/Almaty',
                weekday: 'short',
              }).format(d)
              return (
                <button
                  key={d.toISOString()}
                  type="button"
                  onClick={() => {
                    setTabIndex(i)
                    setSelectedSessionId(null)
                  }}
                  className={`flex min-w-[64px] flex-col items-center rounded-xl border px-3 py-2 ${
                    i === tabIndex
                      ? 'border-accent bg-accentDim text-white'
                      : 'border-border bg-card2 text-muted'
                  }`}
                >
                  <span className="text-lg font-semibold">{dayNum}</span>
                  <span className="text-[10px] uppercase">{weekday}</span>
                </button>
              )
            })}
          </div>

          <ul className="flex flex-col gap-2">
            {sessionsForDay.length === 0 ? (
              <li className="rounded-lg border border-border bg-card2 p-4 text-body text-muted">
                No showtimes on this day.
              </li>
            ) : (
              sessionsForDay.map((s) => (
                <li key={s.id}>
                  <button
                    type="button"
                    onClick={() => onPickSession(s)}
                    className={`flex w-full flex-wrap items-center justify-between gap-3 rounded-lg border px-4 py-3 text-left transition-colors ${
                      selectedSessionId === s.id
                        ? 'border-accent bg-accentDim'
                        : 'border-border bg-card2 hover:border-accent'
                    }`}
                  >
                    <div>
                      <p className="text-body font-medium text-white">
                        {formatTimeLocal(s.start_time)} · Hall {s.hall_name ?? '—'}
                      </p>
                      <p className="mt-0.5 text-[11px] text-muted">
                        {s.available_seats ?? '—'} seats free
                      </p>
                    </div>
                    <span className="text-body font-medium text-accent">
                      {formatPrice(s.price)}
                    </span>
                  </button>
                </li>
              ))
            )}
          </ul>
        </section>
      )}

      <AgeVerifyModal
        open={ageOpen}
        onClose={() => {
          setAgeOpen(false)
          setPendingSession(null)
        }}
        onConfirm={() => {
          setAgeOpen(false)
          if (pendingSession) void proceedToBooking(pendingSession)
          setPendingSession(null)
        }}
      />
    </div>
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
  return (
    <span
      className={`rounded-lg border px-3 py-1 text-[12px] font-medium ${
        active
          ? 'border-accent bg-accent text-white'
          : done
            ? 'border-accent bg-accentDim text-accent'
            : 'border-border bg-card2 text-muted'
      }`}
    >
      Step {n}: {label}
    </span>
  )
}
