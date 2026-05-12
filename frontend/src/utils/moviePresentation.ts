import type { LucideIcon } from 'lucide-react'
import {
  Clapperboard,
  Drama,
  Ghost,
  Laugh,
  Rocket,
  Skull,
  Swords,
} from 'lucide-react'

export function genrePosterClass(genre: string): string {
  const g = genre.toLowerCase()
  if (g.includes('action')) return 'bg-[#1a1010]'
  if (g.includes('horror')) return 'bg-[#120a12]'
  if (g.includes('sci')) return 'bg-[#0a1218]'
  if (g.includes('drama')) return 'bg-[#101012]'
  if (g.includes('anim')) return 'bg-[#10140a]'
  if (g.includes('comedy')) return 'bg-[#12100a]'
  if (g.includes('thrill')) return 'bg-[#140c0c]'
  return 'bg-card2'
}

export function genreIcon(genre: string): LucideIcon {
  const g = genre.toLowerCase()
  if (g.includes('action')) return Swords
  if (g.includes('horror')) return Ghost
  if (g.includes('sci')) return Rocket
  if (g.includes('drama')) return Drama
  if (g.includes('anim')) return Laugh
  if (g.includes('comedy')) return Laugh
  if (g.includes('thrill')) return Skull
  return Clapperboard
}

export const HOME_GENRE_CHIPS: { value: string; label: string }[] = [
  { value: '', label: 'All' },
  { value: 'Action', label: 'Action' },
  { value: 'Drama', label: 'Drama' },
  { value: 'Sci-Fi', label: 'Sci-Fi' },
  { value: 'Horror', label: 'Horror' },
  { value: 'Animation', label: 'Animation' },
  { value: 'Comedy', label: 'Comedy' },
  { value: 'Thriller', label: 'Thriller' },
]

export function movieMatchesGenreFilter(movieGenre: string, chip: string): boolean {
  if (!chip) return true
  return movieGenre.toLowerCase().includes(chip.toLowerCase())
}
