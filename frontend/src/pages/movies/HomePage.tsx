import { useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { MovieGrid } from '../../components/movie/MovieGrid'
import { HomeHero } from '../../components/movie/HomeHero'
import { SectionTitle } from '../../components/ui/SectionTitle'
import { HOME_GENRE_CHIPS, movieMatchesGenreFilter } from '../../utils/moviePresentation'
import {
  getTrendingMovies,
  getUpcomingMovies,
  getNowPlayingMovies,
  getTopRatedMovies,
  getAllCatalogMovies,
} from '../../lib/movieService'
import { filterCatalogByQuery } from '../../lib/movieFilters'
import { buildCatalogPosterMap, catalogToMovie } from '../../lib/catalogAdapter'
import type { CatalogMovie } from '../../types/movie'

function applyCatalogFilters(
  list: CatalogMovie[],
  qRaw: string,
  genre: string,
): CatalogMovie[] {
  const byQ = filterCatalogByQuery(list, qRaw.trim())
  return byQ.filter((c) =>
    movieMatchesGenreFilter(catalogToMovie(c).genre, genre),
  )
}

export function HomePage() {
  const [params] = useSearchParams()
  const qRaw = params.get('q') ?? ''
  const [genre, setGenre] = useState('')

  const posterById = useMemo(
    () => buildCatalogPosterMap(getAllCatalogMovies()),
    [],
  )

  const trending = useMemo(
    () => applyCatalogFilters(getTrendingMovies(), qRaw, genre).map(catalogToMovie),
    [qRaw, genre],
  )
  const upcoming = useMemo(
    () => applyCatalogFilters(getUpcomingMovies(), qRaw, genre).map(catalogToMovie),
    [qRaw, genre],
  )
  const nowPlaying = useMemo(
    () => applyCatalogFilters(getNowPlayingMovies(), qRaw, genre).map(catalogToMovie),
    [qRaw, genre],
  )
  const topRated = useMemo(
    () => applyCatalogFilters(getTopRatedMovies(), qRaw, genre).map(catalogToMovie),
    [qRaw, genre],
  )

  const featuredCatalog = useMemo(
    () => applyCatalogFilters(getTrendingMovies(), qRaw, genre)[0],
    [qRaw, genre],
  )

  return (
    <div className="flex flex-col gap-8">
      <div>
        <h1 className="text-hero font-light text-white">Movies</h1>
        <p className="mt-1 text-body text-muted">
          Browse titles from the local catalog (JSON). Showtimes use your venue API when
          synced.
        </p>
      </div>

      {featuredCatalog ? (
        <HomeHero
          movie={catalogToMovie(featuredCatalog)}
          catalogPosterUrl={featuredCatalog.poster_path}
          catalogBackdropUrl={featuredCatalog.backdrop_path}
        />
      ) : null}

      <div className="flex flex-wrap gap-2">
        {HOME_GENRE_CHIPS.map((chip) => {
          const active = genre === chip.value
          return (
            <button
              key={chip.label}
              type="button"
              onClick={() => setGenre(chip.value)}
              className={`rounded-full border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${
                active
                  ? 'border-accent bg-accent text-white'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              {chip.label}
            </button>
          )
        })}
      </div>

      <section>
        <SectionTitle>Trending</SectionTitle>
        <div className="mt-4">
          {trending.length === 0 ? (
            <p className="text-body text-muted">No movies match.</p>
          ) : (
            <MovieGrid movies={trending} catalogPosterById={posterById} />
          )}
        </div>
      </section>
      <section>
        <SectionTitle>Upcoming</SectionTitle>
        <div className="mt-4">
          {upcoming.length === 0 ? (
            <p className="text-body text-muted">No movies match.</p>
          ) : (
            <MovieGrid movies={upcoming} catalogPosterById={posterById} />
          )}
        </div>
      </section>
      <section>
        <SectionTitle>Now playing</SectionTitle>
        <div className="mt-4">
          {nowPlaying.length === 0 ? (
            <p className="text-body text-muted">No movies match.</p>
          ) : (
            <MovieGrid movies={nowPlaying} catalogPosterById={posterById} />
          )}
        </div>
      </section>
      <section>
        <SectionTitle>Top rated</SectionTitle>
        <div className="mt-4">
          {topRated.length === 0 ? (
            <p className="text-body text-muted">No movies match.</p>
          ) : (
            <MovieGrid movies={topRated} catalogPosterById={posterById} />
          )}
        </div>
      </section>
    </div>
  )
}
