import { useMemo } from 'react'
import type { Seat } from '../../types'
import { useBookingStore } from '../../store/bookingStore'
import { formatPrice } from '../../utils/format'

interface SeatMapProps {
  seats: Seat[]
}

function isVipRow(row: string): boolean {
  const r = row.toUpperCase()
  return r === 'A' || r === 'B'
}

export function SeatMap({ seats }: SeatMapProps) {
  const selectedSeats = useBookingStore((s) => s.selectedSeats)
  const toggleSeat = useBookingStore((s) => s.toggleSeat)
  const session = useBookingStore((s) => s.session)

  const rows = useMemo(() => {
    const byRow = new Map<string, Seat[]>()
    for (const seat of seats) {
      const list = byRow.get(seat.row) ?? []
      list.push(seat)
      byRow.set(seat.row, list)
    }
    const keys = Array.from(byRow.keys()).sort((a, b) =>
      a.localeCompare(b, undefined, { sensitivity: 'base' }),
    )
    return keys.map((row) => ({
      row,
      seats: (byRow.get(row) ?? []).sort((a, b) => a.number - b.number),
    }))
  }, [seats])

  const price = session?.price ?? 0

  function isSelected(seat: Seat): boolean {
    return selectedSeats.some((s) => s.id === seat.id)
  }

  return (
    <div className="flex flex-col items-stretch gap-4">
      <div className="flex w-full flex-col items-center gap-1">
        <div className="h-px w-full bg-border2" aria-hidden />
        <p
          className="text-[11px] uppercase tracking-[0.1em] text-muted"
          style={{ letterSpacing: '0.1em' }}
        >
          screen
        </p>
      </div>
      <div className="flex flex-col gap-2">
        {rows.map(({ row, seats: rowSeats }, rowIdx) => {
          const vip = isVipRow(row)
          const prevVip = rowIdx > 0 ? isVipRow(rows[rowIdx - 1]!.row) : false
          const gapBefore = !vip && prevVip
          return (
            <div
              key={row}
              className={`flex items-center gap-2 ${gapBefore ? 'mt-3' : ''}`}
            >
              <span className="w-4 shrink-0 text-right text-[11px] text-muted">
                {row}
              </span>
              <div className="flex flex-wrap gap-1">
                {rowSeats.map((seat) => {
                  const taken = !seat.is_available
                  const selected = isSelected(seat)
                  const vipSeat = vip
                  let cls =
                    'flex h-[22px] w-[26px] items-center justify-center rounded border text-[9px] font-medium transition-colors duration-150'
                  if (taken && vipSeat) {
                    cls +=
                      ' cursor-not-allowed border-slate-900 bg-[#060d14] text-slate-900 pointer-events-none'
                  } else if (taken) {
                    cls +=
                      ' cursor-not-allowed border-[#2a1515] bg-[#1c1010] text-[#2a1515] pointer-events-none'
                  } else if (selected && vipSeat) {
                    cls += ' border-vip bg-vip text-black'
                  } else if (selected) {
                    cls += ' border-accent bg-accent text-white'
                  } else if (vipSeat) {
                    cls +=
                      ' border-slate-800 bg-[#0a1520] text-slate-500 hover:border-vip hover:text-vip'
                  } else {
                    cls +=
                      ' border-border2 bg-[#1a1a1a] text-muted hover:border-accent hover:text-accent'
                  }
                  return (
                    <button
                      key={seat.id}
                      type="button"
                      disabled={taken}
                      aria-pressed={selected}
                      aria-label={`Row ${row} seat ${seat.number}`}
                      className={cls}
                      onClick={() => toggleSeat(seat)}
                    >
                      {seat.number}
                    </button>
                  )
                })}
              </div>
            </div>
          )
        })}
      </div>
      <div className="flex flex-wrap gap-4 border-t border-border pt-3">
        <LegendItem
          label="Available"
          boxClass="border-border2 bg-[#1a1a1a]"
        />
        <LegendItem label="Selected" boxClass="border-accent bg-accent" />
        <LegendItem
          label="Taken"
          boxClass="border-[#2a1515] bg-[#1c1010]"
        />
        <LegendItem
          label="VIP"
          boxClass="border-slate-800 bg-[#0a1520]"
        />
        <LegendItem label="VIP selected" boxClass="border-vip bg-vip" />
      </div>
      {price > 0 ? (
        <p className="text-[11px] text-muted">Price per seat: {formatPrice(price)}</p>
      ) : null}
    </div>
  )
}

function LegendItem({
  label,
  boxClass,
}: {
  label: string
  boxClass: string
}) {
  return (
    <div className="flex items-center gap-2">
      <span className={`h-2.5 w-3 rounded-sm border ${boxClass}`} aria-hidden />
      <span className="text-[11px] text-muted">{label}</span>
    </div>
  )
}
