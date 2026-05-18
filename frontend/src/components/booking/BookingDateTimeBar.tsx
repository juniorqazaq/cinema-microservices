import { addDays, eachDayOfInterval } from 'date-fns'
import { useMemo } from 'react'
import type { Session } from '../../types'
import { formatPrice, formatTimeLocal, localDayKey, sessionLocalDayKey } from '../../utils/format'

interface BookingDateTimeBarProps {
  sessions: Session[]
  selectedSessionId: string | undefined
  tabIndex: number
  onTabIndexChange: (i: number) => void
  onSelectSession: (s: Session) => void
}

export function BookingDateTimeBar({
  sessions,
  selectedSessionId,
  tabIndex,
  onTabIndexChange,
  onSelectSession,
}: BookingDateTimeBarProps) {
  const days = useMemo(() => {
    const start = new Date()
    start.setHours(0, 0, 0, 0)
    return eachDayOfInterval({ start, end: addDays(start, 6) })
  }, [])

  const selectedDayKey = localDayKey(days[tabIndex] ?? days[0]!)

  const sessionsForDay = useMemo(
    () => sessions.filter((s) => sessionLocalDayKey(s.start_time) === selectedDayKey),
    [sessions, selectedDayKey],
  )

  return (
    <div className="flex flex-col gap-4 rounded-2xl border border-[#1e1e32] bg-[#0D0D1A] p-4">
      <div className="flex gap-2 overflow-x-auto pb-1">
        {days.map((d, i) => {
          const active = i === tabIndex
          const dayNum = new Intl.DateTimeFormat('en', {
            timeZone: 'Asia/Almaty',
            day: 'numeric',
          }).format(d)
          const weekday = new Intl.DateTimeFormat('en', {
            timeZone: 'Asia/Almaty',
            weekday: 'short',
          }).format(d)
          return (
            <button
              key={d.toISOString()}
              type="button"
              onClick={() => onTabIndexChange(i)}
              className={`flex min-w-[64px] flex-col items-center rounded-xl border px-3 py-2 text-center transition-all ${
                active
                  ? 'border-[#7C4DFF] bg-[#1a1030] text-white shadow-[0_0_16px_rgba(124,77,255,0.4)]'
                  : 'border-transparent text-[#9ca3af] hover:text-white'
              }`}
            >
              <span className="text-lg font-semibold">{dayNum}</span>
              <span className="text-[10px] uppercase">{weekday}</span>
            </button>
          )
        })}
      </div>

      <div className="flex flex-wrap gap-2">
        {sessionsForDay.length === 0 ? (
          <p className="text-sm text-[#9ca3af]">No showtimes this day.</p>
        ) : (
          sessionsForDay.map((s) => {
            const active = s.id === selectedSessionId
            return (
              <button
                key={s.id}
                type="button"
                onClick={() => onSelectSession(s)}
                className={`rounded-xl border px-4 py-2 text-sm font-medium transition-all ${
                  active
                    ? 'border-[#7C4DFF] bg-[#1a1030] text-white shadow-[0_0_12px_rgba(124,77,255,0.35)]'
                    : 'border-[#2a2a3d] bg-[#14141f] text-[#9ca3af] hover:border-[#7C4DFF]'
                }`}
              >
                {formatTimeLocal(s.start_time)} · {formatPrice(s.price)}
              </button>
            )
          })
        )}
      </div>
    </div>
  )
}
