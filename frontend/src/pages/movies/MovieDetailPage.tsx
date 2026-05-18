import axios from 'axios'
import { format, parseISO } from 'date-fns'
import { createElement } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useMemo } from 'react'
import { useMovie } from '../../hooks/useMovies'
import { useTmdbMovie } from '../../hooks/useTmdbMovie'
import { BuyTicketButton } from '../../components/movie/BuyTicketButton'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { getErrorMessage } from '../../utils/errorHandler'
import { formatDuration } from '../../utils/format'
import { genreIcon, genrePosterClass } from '../../utils/moviePresentation'
import { tmdbImg } from '../../api/tmdb'
import { getMovieById } from '../../lib/movieService'
import { CatalogMovieDetail } from '../../components/movie/CatalogMovieDetail'

export function MovieDetailPage() {
  const { id } = useParams()
  const catalogMovie = useMemo(() => (id ? getMovieById(id) : undefined), [id])
  const movieQuery = useMovie(id, { enabled: Boolean(id) && !catalogMovie })
  const tmdb = useTmdbMovie(catalogMovie ? undefined : movieQuery.data)

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
  const year = movie.created_at
    ? format(parseISO(movie.created_at), 'yyyy')
    : '—'

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
                    {createElement(genreIcon(movie.genre), {
                      className: 'h-10 w-10 text-white/30',
                      strokeWidth: 1,
                      'aria-hidden': true,
                    })}
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
                <BuyTicketButton />
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  )
}


