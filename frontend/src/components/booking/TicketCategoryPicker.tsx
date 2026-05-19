import {
  TICKET_CATEGORIES,
  ticketPrice,
  type TicketCategoryId,
} from '../../constants/ticketCategories'
import { formatPrice } from '../../utils/format'

interface TicketCategoryPickerProps {
  basePrice: number
  value: TicketCategoryId
  onChange: (id: TicketCategoryId) => void
}

export function TicketCategoryPicker({
  basePrice,
  value,
  onChange,
}: TicketCategoryPickerProps) {
  return (
    <div className="flex flex-col gap-2">
      <p className="text-body font-medium text-white">Ticket category</p>
      <div className="grid gap-2 sm:grid-cols-2">
        {TICKET_CATEGORIES.map((cat) => {
          const price = ticketPrice(basePrice, cat.id)
          const active = value === cat.id
          return (
            <button
              key={cat.id}
              type="button"
              onClick={() => onChange(cat.id)}
              className={`rounded-lg border px-3 py-2 text-left text-body transition-colors ${
                active
                  ? 'border-accent bg-accentDim text-white'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              <span className="font-medium">{cat.label}</span>
              <span className="mt-0.5 block text-[11px] text-muted">
                {cat.description} · {formatPrice(price)}
              </span>
            </button>
          )
        })}
      </div>
    </div>
  )
}
