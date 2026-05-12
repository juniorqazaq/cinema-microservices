/** Rich catalog movie (local JSON / Netflix-style metadata). */
export interface CatalogReview {
  author: string
  rating: number
  content: string
}

export interface CatalogCastMember {
  name: string
  character: string
  profile_path: string
}

export interface CatalogMovie {
  id: number
  title: string
  original_title: string
  overview: string
  /** TMDB poster file path — use `tmdbImage(poster_path, 'w500')` etc. */
  poster_path: string
  /** TMDB backdrop file path — use `tmdbImage(backdrop_path, 'original')` etc. */
  backdrop_path: string
  release_date: string
  vote_average: number
  vote_count: number
  popularity: number
  genre_ids: number[]
  adult: boolean
  original_language: string
  runtime: number
  status: string
  tagline: string
  reviews: CatalogReview[]
  cast: CatalogCastMember[]
  trailer: string
}

export type CatalogCategory = 'trending' | 'upcoming' | 'nowPlaying' | 'topRated'

export type SortCatalogBy = 'vote_average' | 'popularity' | 'release_date' | 'title'
