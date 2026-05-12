import { Link } from 'react-router-dom'
import { createElement } from 'react'
import type { Movie } from '../../types'
import { formatDuration } from '../../utils/format'
import { genreIcon, genrePosterClass } from '../../utils/moviePresentation'
import { useTmdbMovie } from '../../hooks/useTmdbMovie'
import { tmdbImg } from '../../api/tmdb'
import { posterUrl } from '../../lib/movieImages'

interface MovieCardProps {
  movie: Movie
  /** Full poster URL (e.g. from local JSON catalog). Skips TMDB. */
  catalogPosterUrl?: string | null
}

export function MovieCard({ movie, catalogPosterUrl }: MovieCardProps) {
  const tmdb = useTmdbMovie(catalogPosterUrl ? undefined : movie)
  const poster = catalogPosterUrl
    ? posterUrl({ poster_path: catalogPosterUrl }, 'w500')
    : tmdbImg(tmdb.data?.poster_path ?? null, 'w500')

  return (
    <Link
      to={`/movies/${movie.id}`}
      className="flex flex-col overflow-hidden rounded-xl border border-border bg-card transition-colors duration-150 hover:border-accent"
    >
      <div className="relative flex h-[200px] w-full items-center justify-center overflow-hidden bg-card2">
        {poster ? (
          <>
            <img
              src={poster}
              alt=""
              className="absolute inset-0 h-full w-full object-cover"
              loading="lazy"
            />
            <div className="absolute inset-0 bg-black/25" aria-hidden />
          </>
        ) : (
          <div
            className={`absolute inset-0 flex items-center justify-center ${genrePosterClass(movie.genre)}`}
          >
            {createElement(genreIcon(movie.genre), {
              className: 'h-12 w-12 text-white/25',
              strokeWidth: 1,
              'aria-hidden': true,
            })}
          </div>
        )}
        <div className="absolute right-2 top-2 rounded-md border border-accent bg-accent px-2 py-0.5 text-[11px] font-medium text-white">
          <span className="text-pink-400" aria-hidden>
            ★
          </span>{' '}
          {movie.rating.toFixed(1)}
        </div>
      </div>
      <div className="flex flex-col gap-1 p-3">
        <h2 className="line-clamp-2 text-card-title font-medium text-white">
          {movie.title}
        </h2>
        <p className="text-body text-muted">
          {movie.genre} · {formatDuration(movie.duration)}
        </p>
      </div>
    </Link>
  )
}
