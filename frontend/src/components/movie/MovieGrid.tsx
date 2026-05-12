import type { Movie } from '../../types'
import { MovieCard } from './MovieCard'

export function MovieGridSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: 8 }).map((_, i) => (
        <div
          key={i}
          className="overflow-hidden rounded-lg border border-border bg-card"
          aria-hidden="true"
        >
          <div className="h-[200px] animate-pulse bg-card2" />
          <div className="space-y-2 p-3">
            <div className="h-4 w-3/4 animate-pulse rounded bg-card2" />
            <div className="h-3 w-1/2 animate-pulse rounded bg-card2" />
          </div>
        </div>
      ))}
    </div>
  )
}

export function MovieGrid({
  movies,
  catalogPosterById,
}: {
  movies: Movie[]
  catalogPosterById?: Record<string, string | undefined>
}) {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
      {movies.map((m) => (
        <MovieCard
          key={m.id}
          movie={m}
          catalogPosterUrl={catalogPosterById?.[m.id]}
        />
      ))}
    </div>
  )
}
