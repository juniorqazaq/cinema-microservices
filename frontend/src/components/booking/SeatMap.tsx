import { useMemo } from 'react'
import type { Seat } from '../../types'
import { useBookingStore } from '../../store/bookingStore'
import { formatPrice } from '../../utils/format'
import { ticketPrice } from '../../constants/ticketCategories'

interface SeatMapProps {
  seats: Seat[]
}

function isVipRow(row: string): boolean {
  const r = row.toUpperCase()
  return r === 'A' || r === 'B'
}

function rowSortKey(row: string): number {
  const r = row.toUpperCase()
  if (r.length === 1 && r >= 'A' && r <= 'Z') return r.charCodeAt(0)
  const n = parseInt(r.replace(/\D/g, ''), 10)
  return Number.isFinite(n) ? 100 + n : 999
}

export function SeatMap({ seats }: SeatMapProps) {
  const selectedSeats = useBookingStore((s) => s.selectedSeats)
  const toggleSeat = useBookingStore((s) => s.toggleSeat)
  const session = useBookingStore((s) => s.session)
  const ticketCategory = useBookingStore((s) => s.ticketCategory)

  const rows = useMemo(() => {
    const byRow = new Map<string, Seat[]>()
    for (const seat of seats) {
      const list = byRow.get(seat.row) ?? []
      list.push(seat)
      byRow.set(seat.row, list)
    }
    const keys = Array.from(byRow.keys()).sort(
      (a, b) => rowSortKey(a) - rowSortKey(b),
    )
    return keys.map((row) => ({
      row,
      seats: (byRow.get(row) ?? []).sort((a, b) => a.number - b.number),
    }))
  }, [seats])

  const basePrice = session?.price ?? 0
  const unitPrice = ticketPrice(basePrice, ticketCategory)
  const total = selectedSeats.length * unitPrice

  function isSelected(seat: Seat): boolean {
    return selectedSeats.some((s) => s.id === seat.id)
  }

  return (
    <div className="flex flex-col gap-5 rounded-2xl border border-[#1e1e32] bg-[#0D0D1A] p-5">
      <div className="flex flex-wrap items-center justify-center gap-6 text-[11px] text-[#9ca3af]">
        <LegendDot color="#00E5FF" label="Selected" />
        <LegendDot color="#7C4DFF" label="Reserved" />
        <LegendDot color="#2A2A3D" label="Available" />
        <LegendDot color="#d4a017" label="VIP" />
      </div>

      <p className="text-center text-lg font-medium tracking-wide text-white">
        How Many Seats?
      </p>

      <div
        className="relative mx-auto w-full max-w-2xl"
        style={{ perspective: '700px' }}
      >
        <div
          className="mx-auto mb-6 h-1 w-3/4 rounded-full"
          style={{
            background:
              'linear-gradient(90deg, transparent, #3b82f6, #60a5fa, #3b82f6, transparent)',
            boxShadow: '0 0 24px rgba(59, 130, 246, 0.8)',
          }}
        />
        <p className="mb-8 text-center text-[10px] uppercase tracking-[0.2em] text-[#60a5fa]">
          Screen
        </p>

        <div
          className="flex flex-col items-center gap-1.5"
          style={{ transform: 'rotateX(22deg)', transformOrigin: 'center top' }}
        >
          {rows.map(({ row, seats: rowSeats }, rowIdx) => {
            const vip = isVipRow(row)
            const arc = Math.sin((rowIdx / Math.max(rows.length - 1, 1)) * Math.PI) * 12
            return (
              <div
                key={row}
                className="flex items-center justify-center gap-1"
                style={{ transform: `translateY(${rowIdx * 2}px)` }}
              >
                <span className="mr-2 w-4 shrink-0 text-right text-[10px] text-[#6b7280]">
                  {row}
                </span>
                {rowSeats.map((seat, seatIdx) => {
                  const taken = !seat.is_available
                  const selected = isSelected(seat)
                  const offset =
                    (seatIdx - (rowSeats.length - 1) / 2) * (1.2 + arc * 0.08)
                  let bg = '#2A2A3D'
                  let border = '#3d3d52'
                  let glow = 'none'
                  if (taken) {
                    bg = '#7C4DFF'
                    border = '#9d6fff'
                    glow = '0 0 8px rgba(124, 77, 255, 0.6)'
                  } else if (selected) {
                    bg = '#00E5FF'
                    border = '#00E5FF'
                    glow = '0 0 10px rgba(0, 229, 255, 0.7)'
                  } else if (vip) {
                    bg = '#1a1520'
                    border = '#d4a017'
                  }
                  return (
                    <button
                      key={seat.id}
                      type="button"
                      disabled={taken}
                      aria-pressed={selected}
                      aria-label={`Row ${row} seat ${seat.number}`}
                      className="h-5 w-5 shrink-0 rounded-sm text-[8px] font-semibold transition-transform duration-150 hover:scale-110 disabled:cursor-not-allowed disabled:opacity-90"
                      style={{
                        backgroundColor: bg,
                        border: `1px solid ${border}`,
                        boxShadow: glow,
                        color: selected ? '#0D0D1A' : taken ? '#fff' : '#9ca3af',
                        transform: `translateX(${offset}px)`,
                      }}
                      onClick={() => toggleSeat(seat)}
                    >
                      {seat.number}
                    </button>
                  )
                })}
              </div>
            )
          })}
        </div>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-[#1e1e32] pt-4">
        <p className="text-sm text-[#9ca3af]">
          {selectedSeats.length === 0
            ? 'Select seats on the map'
            : `${selectedSeats.length} seat(s): ${selectedSeats.map((s) => `${s.row}${s.number}`).join(', ')}`}
        </p>
        {basePrice > 0 ? (
          <p className="text-sm font-medium text-white">
            Total: {formatPrice(total || unitPrice)}
          </p>
        ) : null}
      </div>
    </div>
  )
}

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <span className="inline-flex items-center gap-1.5">
      <span
        className="h-2.5 w-2.5 rounded-full"
        style={{ backgroundColor: color, boxShadow: `0 0 6px ${color}` }}
        aria-hidden
      />
      {label}
    </span>
  )
}
