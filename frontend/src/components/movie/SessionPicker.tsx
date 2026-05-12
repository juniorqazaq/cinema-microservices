import type { Session } from '../../types'
import { formatPrice, formatTime } from '../../utils/format'

interface SessionPickerProps {
  sessions: Session[]
  hallLabel: (hallId: string) => string
  onSelect: (session: Session) => void
}

export function SessionPicker({
  sessions,
  hallLabel,
  onSelect,
}: SessionPickerProps) {
  if (sessions.length === 0) {
    return (
      <p className="text-body text-muted">No sessions for this date.</p>
    )
  }

  return (
    <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
      {sessions.map((s) => (
        <button
          key={s.id}
          type="button"
          className="flex flex-col items-start gap-1 rounded-lg border border-border bg-card2 px-3 py-2 text-left transition-colors duration-150 hover:border-accent"
          onClick={() => onSelect(s)}
        >
          <span className="text-[16px] font-medium text-white">
            {formatTime(s.start_time)}
          </span>
          <span className="text-[11px] text-muted">
            Hall {hallLabel(s.hall_id)}
          </span>
          <span className="text-body text-accent">{formatPrice(s.price)}</span>
        </button>
      ))}
    </div>
  )
}
