import { Link } from 'react-router-dom'
import { format, parseISO } from 'date-fns'
import type { CatalogMovie } from '../../types/movie'
import { backdropUrl, posterUrl } from '../../lib/movieImages'
import { genreIdsToLabel } from '../../lib/catalogAdapter'

const buyTicketBtnClass =
  'inline-flex min-w-[140px] items-center justify-center rounded-lg border border-accent bg-accent px-4 py-2.5 text-body font-semibold text-white shadow-sm transition-colors hover:border-accentHover hover:bg-accentHover'

export function CatalogMovieDetail({ movie }: { movie: CatalogMovie }) {
  const year = movie.release_date
    ? format(parseISO(`${movie.release_date}T12:00:00`), 'yyyy')
    : '—'

  const trailerYoutubeId = (() => {
    if (!movie.trailer) return ''
    try {
      return new URL(movie.trailer).searchParams.get('v') ?? ''
    } catch {
      return ''
    }
  })()

  return (
    <div className="flex flex-col gap-8">
      <Link
        to="/movies"
        className="inline-flex text-body text-muted hover:text-accent"
      >
        ← Back to movies
      </Link>

      <div className="relative overflow-hidden rounded-xl border border-border">
        <div className="relative min-h-[200px] md:min-h-[280px]">
          <img
            src={backdropUrl(movie, 'original')}
            alt=""
            className="absolute inset-0 h-full w-full object-cover"
            loading="eager"
          />
          <div className="absolute inset-0 bg-black/65" aria-hidden />
          <div className="relative flex flex-col gap-4 p-6 md:flex-row md:items-end">
            <div className="flex shrink-0 justify-center md:justify-start">
              <div className="relative h-[200px] w-[133px] overflow-hidden rounded-lg border border-border bg-card2 md:h-[240px] md:w-[160px]">
                <img
                  src={posterUrl(movie, 'w500')}
                  alt=""
                  className="h-full w-full object-cover"
                />
              </div>
            </div>
            <div className="min-w-0 flex-1 pb-1">
              <h1 className="text-hero font-light text-white md:text-[36px]">{movie.title}</h1>
              <p className="mt-2 text-body text-muted">{movie.tagline}</p>
              <div className="mt-3 flex flex-wrap items-center gap-2">
                <span className="rounded-md border border-accent bg-accent px-2 py-0.5 text-[11px] font-medium text-white">
                  <span className="text-pink-400" aria-hidden>
                    ★
                  </span>{' '}
                  {movie.vote_average.toFixed(1)}
                </span>
                <p className="text-body text-white/90">
                  {genreIdsToLabel(movie.genre_ids)} ·{' '}
                  {movie.runtime ? `${movie.runtime} min` : 'TBA'} · {year} ·{' '}
                  {movie.vote_count.toLocaleString()} votes
                </p>
              </div>
              <p className="mt-3 text-body text-white/85">{movie.overview}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                <a href="#catalog-tickets-info" className={buyTicketBtnClass}>
                  Buy ticket
                </a>
                {movie.trailer ? (
                  <a
                    href={movie.trailer}
                    target="_blank"
                    rel="noreferrer noopener"
                    className="inline-flex items-center justify-center rounded-lg border border-border2 bg-card2/90 px-4 py-2.5 text-body font-medium text-white transition-colors hover:border-accent"
                  >
                    Watch trailer
                  </a>
                ) : null}
              </div>
            </div>
          </div>
        </div>
      </div>

      <section className="rounded-lg border border-border bg-card p-4">
        <h2 className="text-section font-light text-white">Cast</h2>
        <div className="mt-4 flex gap-4 overflow-x-auto pb-2">
          {movie.cast.map((c) => (
            <div key={`${c.name}-${c.character}`} className="w-24 shrink-0 text-center">
              <img
                src={c.profile_path}
                alt=""
                className="mx-auto h-20 w-20 rounded-full border border-border object-cover"
                loading="lazy"
              />
              <p className="mt-2 line-clamp-2 text-[11px] font-medium text-white">{c.name}</p>
              <p className="line-clamp-2 text-[10px] text-muted">{c.character}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="rounded-lg border border-border bg-card p-4">
        <h2 className="text-section font-light text-white">Reviews</h2>
        <ul className="mt-4 flex flex-col gap-3">
          {movie.reviews.map((r, i) => (
            <li
              key={`${r.author}-${i}`}
              className="rounded-lg border border-border2 bg-card2/60 p-3"
            >
              <div className="flex items-center justify-between gap-2">
                <p className="text-body font-medium text-white">{r.author}</p>
                <span className="text-body text-accent">{r.rating.toFixed(1)}/10</span>
              </div>
              <p className="mt-2 text-body text-muted">{r.content}</p>
            </li>
          ))}
        </ul>
      </section>

      {trailerYoutubeId ? (
        <section className="overflow-hidden rounded-xl border border-border bg-card">
          <div className="aspect-video w-full max-w-3xl">
            <iframe
              title="Trailer"
              className="h-full w-full"
              src={`https://www.youtube-nocookie.com/embed/${trailerYoutubeId}`}
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            />
          </div>
        </section>
      ) : null}

      <section
        id="catalog-tickets-info"
        className="scroll-mt-28 rounded-lg border border-border bg-card2/50 p-4 text-body text-muted"
      >
        <p>
          Showtimes use the venue API. When this title is synced to your cinema backend,
          sessions will appear on this page.
        </p>
      </section>
    </div>
  )
}
