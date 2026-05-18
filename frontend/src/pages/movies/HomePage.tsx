import { useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import type { SortCatalogBy } from '../../types/movie'
import { sortCatalogMovies } from '../../lib/movieFilters'
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
  sortBy: SortCatalogBy,
): CatalogMovie[] {
  const byQ = filterCatalogByQuery(list, qRaw.trim())
  const byGenre = byQ.filter((c) =>
    movieMatchesGenreFilter(catalogToMovie(c).genre, genre),
  )
  return sortCatalogMovies(byGenre, sortBy, 'desc')
}

export function HomePage() {
  const [params] = useSearchParams()
  const qRaw = params.get('q') ?? ''
  const [genre, setGenre] = useState('')
  const [sortBy, setSortBy] = useState<SortCatalogBy>('vote_average')

  const posterById = useMemo(
    () => buildCatalogPosterMap(getAllCatalogMovies()),
    [],
  )

  const trending = useMemo(
    () => applyCatalogFilters(getTrendingMovies(), qRaw, genre, sortBy).map(catalogToMovie),
    [qRaw, genre, sortBy],
  )
  const upcoming = useMemo(
    () => applyCatalogFilters(getUpcomingMovies(), qRaw, genre, sortBy).map(catalogToMovie),
    [qRaw, genre, sortBy],
  )
  const nowPlaying = useMemo(
    () => applyCatalogFilters(getNowPlayingMovies(), qRaw, genre, sortBy).map(catalogToMovie),
    [qRaw, genre, sortBy],
  )
  const topRated = useMemo(
    () => applyCatalogFilters(getTopRatedMovies(), qRaw, genre, sortBy).map(catalogToMovie),
    [qRaw, genre, sortBy],
  )

  const featuredCatalog = useMemo(
    () => applyCatalogFilters(getTrendingMovies(), qRaw, genre, sortBy)[0],
    [qRaw, genre, sortBy],
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

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
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
        <label className="flex items-center gap-2 text-body text-muted">
          Sort by rating
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as SortCatalogBy)}
            className="rounded-lg border border-border bg-card2 px-3 py-1.5 text-white"
          >
            <option value="vote_average">Rating</option>
            <option value="popularity">Popularity</option>
            <option value="release_date">Release date</option>
            <option value="title">Title</option>
          </select>
        </label>
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
