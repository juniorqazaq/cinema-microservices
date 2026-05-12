import axios from 'axios'
import { addDays, eachDayOfInterval, format, parseISO } from 'date-fns'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useMemo, useState } from 'react'
import { useQueries } from '@tanstack/react-query'
import { useMovie } from '../../hooks/useMovies'
import { useSessions } from '../../hooks/useSessions'
import { useTmdbMovie } from '../../hooks/useTmdbMovie'
import { SessionPicker } from '../../components/movie/SessionPicker'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import { getErrorMessage, mapApiError } from '../../utils/errorHandler'
import { formatDuration } from '../../utils/format'
import { genreIcon, genrePosterClass } from '../../utils/moviePresentation'
import { fetchHall } from '../../api/movies'
import { tmdbImg } from '../../api/tmdb'
import { useBookingStore } from '../../store/bookingStore'
import { getMovieById } from '../../lib/movieService'
import { CatalogMovieDetail } from '../../components/movie/CatalogMovieDetail'

const buyTicketBtnClass =
  'inline-flex min-w-[140px] items-center justify-center rounded-lg border border-accent bg-accent px-4 py-2.5 text-body font-semibold text-white shadow-sm transition-colors hover:border-accentHover hover:bg-accentHover'

export function MovieDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const setMovie = useBookingStore((s) => s.setMovie)
  const setSession = useBookingStore((s) => s.setSession)
  const setHall = useBookingStore((s) => s.setHall)
  const setStep = useBookingStore((s) => s.setStep)

  const catalogMovie = useMemo(() => (id ? getMovieById(id) : undefined), [id])
  const movieQuery = useMovie(id, { enabled: Boolean(id) && !catalogMovie })
  const tmdb = useTmdbMovie(catalogMovie ? undefined : movieQuery.data)
  const [tabIndex, setTabIndex] = useState(0)

  const days = useMemo(() => {
    const start = new Date()
    start.setHours(0, 0, 0, 0)
    return eachDayOfInterval({ start, end: addDays(start, 6) })
  }, [])

  const selectedDate = days[tabIndex] ?? days[0]
  const dateParam = format(selectedDate, 'yyyy-MM-dd')

  const sessionsQuery = useSessions(
    { movie_id: id, date: dateParam },
    {
      enabled:
        Boolean(id) &&
        !catalogMovie &&
        Boolean(movieQuery.data) &&
        Boolean(dateParam),
    },
  )

  const sessions = sessionsQuery.data ?? []
  const uniqueHallIds = useMemo(
    () => Array.from(new Set(sessions.map((s) => s.hall_id))),
    [sessions],
  )

  const hallQueries = useQueries({
    queries: uniqueHallIds.map((hallId) => ({
      queryKey: ['hall', hallId],
      queryFn: () => fetchHall(hallId),
      enabled: Boolean(hallId) && sessions.length > 0,
    })),
  })

  const hallsById = useMemo(() => {
    const map: Record<string, { name: string } | undefined> = {}
    uniqueHallIds.forEach((hid, i) => {
      map[hid] = hallQueries[i]?.data
    })
    return map
  }, [uniqueHallIds, hallQueries])

  const notFound =
    !catalogMovie &&
    movieQuery.isError &&
    axios.isAxiosError(movieQuery.error) &&
    movieQuery.error.response?.status === 404

  if (catalogMovie) {
    return <CatalogMovieDetail movie={catalogMovie} />
  }

  if (notFound) {
    return (
      <div className="rounded-lg border border-border bg-card p-6 text-body text-muted">
        Not found
      </div>
    )
  }

  if (movieQuery.isError) {
    return (
      <div className="rounded-lg border border-border bg-card p-6">
        <ErrorBanner message={getErrorMessage(movieQuery.error)} />
      </div>
    )
  }

  if (!catalogMovie && (movieQuery.isLoading || !movieQuery.data)) {
    return (
      <div className="flex flex-col gap-4">
        <div className="h-8 w-2/3 animate-pulse rounded-lg bg-card2" />
        <div className="h-40 animate-pulse rounded-lg bg-card2" />
      </div>
    )
  }

  const movie = movieQuery.data!
  const backdrop = tmdbImg(tmdb.data?.backdrop_path ?? null, 'w780')
  const posterTmdb = tmdbImg(tmdb.data?.poster_path ?? null, 'w500')
  const Icon = genreIcon(movie.genre)
  const year = movie.created_at
    ? format(parseISO(movie.created_at), 'yyyy')
    : '—'

  async function handleSessionPick(session: (typeof sessions)[0]) {
    setMovie(movie)
    setSession(session)
    const hall = await fetchHall(session.hall_id)
    setHall(hall)
    setStep(3)
    navigate('/booking')
  }

  return (
    <div className="flex flex-col gap-6">
      <Link
        to="/movies"
        className="inline-flex text-body text-muted hover:text-accent"
      >
        ← Back to movies
      </Link>

      <div className="relative overflow-hidden rounded-xl border border-border">
        <div className="relative min-h-[200px] md:min-h-[260px]">
          {backdrop ? (
            <img
              src={backdrop}
              alt=""
              className="absolute inset-0 h-full w-full object-cover"
            />
          ) : (
            <div
              className={`absolute inset-0 ${genrePosterClass(movie.genre)}`}
              aria-hidden
            />
          )}
          <div className="absolute inset-0 bg-black/65" aria-hidden />
          <div className="relative flex flex-col gap-4 p-6 md:flex-row md:items-end">
            <div className="flex shrink-0 justify-center md:justify-start">
              <div className="relative h-[200px] w-[133px] overflow-hidden rounded-lg border border-border bg-card2 md:h-[240px] md:w-[160px]">
                {posterTmdb ? (
                  <img
                    src={posterTmdb}
                    alt=""
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div
                    className={`flex h-full w-full items-center justify-center ${genrePosterClass(movie.genre)}`}
                  >
                    <Icon
                      className="h-10 w-10 text-white/30"
                      strokeWidth={1}
                      aria-hidden
                    />
                  </div>
                )}
              </div>
            </div>
            <div className="min-w-0 flex-1 pb-1">
              <h1 className="text-hero font-light text-white md:text-[36px]">
                {movie.title}
              </h1>
              <div className="mt-3 flex flex-wrap items-center gap-2">
                <span className="rounded-md border border-accent bg-accent px-2 py-0.5 text-[11px] font-medium text-white">
                  <span className="text-pink-400" aria-hidden>
                    ★
                  </span>{' '}
                  {movie.rating.toFixed(1)}
                </span>
                <p className="text-body text-white/90">
                  {movie.genre} · {formatDuration(movie.duration)} · {year}
                </p>
              </div>
              <p className="mt-3 line-clamp-3 text-body text-white/80 md:line-clamp-none">
                {movie.description}
              </p>
              <div className="mt-4 flex flex-wrap gap-2">
                <a href="#movie-sessions" className={buyTicketBtnClass}>
                  Buy ticket
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div id="movie-sessions" className="min-w-0 scroll-mt-28">
        <div className="my-2 border-t border-border md:hidden" />
        <h2 className="text-section font-light text-white">Sessions</h2>
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
          {sessionsQuery.isLoading ? (
            <Spinner className="h-6 w-6" />
          ) : sessionsQuery.isError ? (
            <p className="text-body text-danger">
              {mapApiError(
                axios.isAxiosError(sessionsQuery.error)
                  ? sessionsQuery.error.response?.status
                  : undefined,
              ) ?? 'Could not load sessions.'}
            </p>
          ) : (
            <SessionPicker
              sessions={sessions}
              hallLabel={(hid) => hallsById[hid]?.name ?? hid.slice(0, 6)}
              onSelect={(s) => void handleSessionPick(s)}
            />
          )}
        </div>
      </div>
    </div>
  )
}
