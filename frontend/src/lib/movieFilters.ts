import type { CatalogMovie, SortCatalogBy } from '../types/movie'

function norm(s: string): string {
  return s.trim().toLowerCase()
}

export function filterCatalogByQuery(
  movies: CatalogMovie[],
  query: string,
): CatalogMovie[] {
  const q = norm(query)
  if (!q) return movies
  return movies.filter((m) => {
    const hay = [
      m.title,
      m.original_title,
      m.overview,
      m.tagline,
      ...m.cast.map((c) => `${c.name} ${c.character}`),
    ]
      .join(' ')
      .toLowerCase()
    return hay.includes(q)
  })
}

export function sortCatalogMovies(
  movies: CatalogMovie[],
  by: SortCatalogBy,
  order: 'asc' | 'desc' = 'desc',
): CatalogMovie[] {
  const mul = order === 'desc' ? -1 : 1
  const copy = [...movies]
  copy.sort((a, b) => {
    switch (by) {
      case 'vote_average':
        return (a.vote_average - b.vote_average) * mul
      case 'popularity':
        return (a.popularity - b.popularity) * mul
      case 'title':
        return a.title.localeCompare(b.title) * mul
      case 'release_date':
        return a.release_date.localeCompare(b.release_date) * mul
      default:
        return 0
    }
  })
  return copy
}
