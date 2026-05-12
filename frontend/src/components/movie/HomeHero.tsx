import { Link } from 'react-router-dom'
import { formatDuration } from '../../utils/format'
import { useTmdbMovie } from '../../hooks/useTmdbMovie'
import { tmdbImg } from '../../api/tmdb'
import { posterUrl, backdropUrl } from '../../lib/movieImages'
import type { Movie } from '../../types'

const btnPrimary =
  'inline-flex min-w-[120px] items-center justify-center gap-2 rounded-lg border border-accent bg-accent px-4 py-2 text-body font-medium text-white transition-colors duration-150 hover:border-accentHover hover:bg-accentHover'
const btnGhost =
  'inline-flex items-center justify-center gap-2 rounded-lg border border-border2 bg-card2/80 px-4 py-2 text-body font-medium text-white transition-colors duration-150 hover:border-accent'

interface HomeHeroProps {
  movie: Movie
  /** Full poster URL from local catalog (skips TMDB). */
  catalogPosterUrl?: string | null
  catalogBackdropUrl?: string | null
}

export function HomeHero({
  movie,
  catalogPosterUrl,
  catalogBackdropUrl,
}: HomeHeroProps) {
  const tmdb = useTmdbMovie(catalogPosterUrl ? undefined : movie)
  const tmdbBackdrop = tmdbImg(tmdb.data?.backdrop_path ?? null, 'w780')
  const backdropSrc = catalogBackdropUrl
    ? backdropUrl({ backdrop_path: catalogBackdropUrl }, 'original')
    : tmdbBackdrop
  const poster =
    catalogPosterUrl ||
    tmdbImg(tmdb.data?.poster_path ?? null, 'w342') ||
    null

  return (
    <section className="relative overflow-hidden rounded-xl border border-border">
      <div className="relative min-h-[220px] md:min-h-[280px]">
        {backdropSrc ? (
          <img
            src={backdropSrc}
            alt=""
            className="absolute inset-0 h-full w-full object-cover"
            loading="eager"
          />
        ) : (
          <div className="absolute inset-0 bg-card2" aria-hidden />
        )}
        <div className="absolute inset-0 bg-black/70" aria-hidden />
        <div className="relative flex flex-col gap-4 p-6 md:flex-row md:items-end md:justify-between">
          <div className="flex max-w-2xl flex-col gap-3 md:flex-row md:items-end md:gap-6">
            {poster ? (
              <img
                src={
                  catalogPosterUrl
                    ? posterUrl({ poster_path: catalogPosterUrl }, 'w500')
                    : poster
                }
                alt=""
                className="hidden h-40 w-[107px] shrink-0 rounded-lg border border-border object-cover md:block"
              />
            ) : null}
            <div>
              <p className="text-[11px] font-medium uppercase tracking-wide text-muted">
                Featured
              </p>
              <h2 className="mt-1 text-hero font-light text-white md:text-[36px]">
                {movie.title}
              </h2>
              <p className="mt-2 text-body text-muted">
                {movie.genre} · {formatDuration(movie.duration)} ·{' '}
                <span className="text-pink-400" aria-hidden>
                  ★
                </span>{' '}
                {movie.rating.toFixed(1)}
              </p>
              <p className="mt-2 line-clamp-2 text-body text-muted">
                {movie.description}
              </p>
            </div>
          </div>
          <div className="flex shrink-0 flex-wrap gap-2">
            <Link to={`/movies/${movie.id}`} className={btnPrimary}>
              Book now
            </Link>
            <Link to={`/movies/${movie.id}`} className={btnGhost}>
              Sessions
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}
