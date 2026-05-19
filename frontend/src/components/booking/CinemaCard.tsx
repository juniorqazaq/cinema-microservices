import { cinemaAddress } from '../../constants/cinemas'

export interface CinemaInfo {
  name: string
  city: string
  hallCount: number
}

interface CinemaCardProps {
  cinema: CinemaInfo
  selected: boolean
  onSelect: () => void
}

export function CinemaCard({ cinema, selected, onSelect }: CinemaCardProps) {
  return (
    <button
      type="button"
      onClick={onSelect}
      className={`flex w-full flex-col rounded-xl border p-4 text-left transition-all ${
        selected
          ? 'border-accent bg-accentDim shadow-[0_0_20px_rgba(124,77,255,0.25)]'
          : 'border-border bg-card2 hover:border-accent'
      }`}
    >
      <p className="text-body font-semibold text-white">{cinema.name}</p>
      <p className="mt-1 text-[12px] text-muted">
        {cinema.hallCount} hall{cinema.hallCount === 1 ? '' : 's'}
      </p>
      <p className="mt-2 text-[11px] text-muted">
        {cinemaAddress(cinema.name, cinema.city)}
      </p>
    </button>
  )
}
