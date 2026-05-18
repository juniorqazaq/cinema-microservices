import { useState } from 'react'
import axios from 'axios'
import { useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '../../store/authStore'
import { useBookingStore } from '../../store/bookingStore'
import type { Hall, Movie, Seat, Session } from '../../types'
import { ticketPrice, type TicketCategoryId } from '../../constants/ticketCategories'
import { formatPrice } from '../../utils/format'
import { Button } from '../ui/Button'
import { Spinner } from '../ui/Spinner'
import { createBooking, confirmPayment } from '../../api/bookings'
import { fetchProfile, topUpWallet } from '../../api/wallet'
import { TopUpModal } from './TopUpModal'
import { TicketCategoryPicker } from './TicketCategoryPicker'

type PayTab = 'wallet' | 'card' | 'kaspi'

interface PaymentPanelProps {
  movie: Movie
  session: Session
  hall: Hall
  selectedSeats: Seat[]
  onPaid: (info: { bookingId: string; totalPaid: number }) => void
  onError: (msg: string) => void
}

export function PaymentPanel({
  movie,
  session,
  hall,
  selectedSeats,
  onPaid,
  onError,
}: PaymentPanelProps) {
  const user = useAuthStore((s) => s.user)
  const setUser = useAuthStore((s) => s.setUser)
  const ticketCategory = useBookingStore((s) => s.ticketCategory)
  const setTicketCategory = useBookingStore((s) => s.setTicketCategory)
  const [tab, setTab] = useState<PayTab>('wallet')
  const [loading, setLoading] = useState(false)
  const [topUpOpen, setTopUpOpen] = useState(false)
  const queryClient = useQueryClient()

  const unitPrice = ticketPrice(session.price, ticketCategory)
  const total = selectedSeats.length * unitPrice
  const balance = user?.balance ?? 0
  const sufficient = balance >= total

  async function handleConfirm() {
    if (tab === 'wallet' && !sufficient) {
      onError('Insufficient balance. Top up your wallet to continue.')
      return
    }

    setLoading(true)
    try {
      let lastBookingId = ''
      for (const seat of selectedSeats) {
        const booking = await createBooking(
          session.id,
          seat.id,
          ticketCategory,
          unitPrice,
        )
        lastBookingId = booking.id
        await confirmPayment(booking.id, unitPrice, tab === 'wallet')
      }
      if (tab === 'wallet') {
        try {
          const profile = await fetchProfile()
          setUser(profile)
        } catch {
          if (user) setUser({ ...user, balance: Math.max(0, balance - total) })
        }
      }
      await queryClient.invalidateQueries({ queryKey: ['bookings'] })
      onPaid({ bookingId: lastBookingId, totalPaid: total })
    } catch (e) {
      if (axios.isAxiosError(e) && e.response?.status === 409) {
        onError('Seat was just taken. Please go back and choose another.')
      } else if (axios.isAxiosError(e)) {
        const msg =
          (e.response?.data as { error?: string } | undefined)?.error ??
          'Payment failed.'
        onError(msg)
      } else {
        onError('Something went wrong. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  const tabs: { id: PayTab; label: string }[] = [
    { id: 'wallet', label: 'Wallet' },
    { id: 'card', label: 'Card' },
    { id: 'kaspi', label: 'Kaspi QR' },
  ]

  const payDisabled =
    loading || selectedSeats.length === 0 || (tab === 'wallet' && !sufficient)

  return (
    <>
      <div className="flex max-w-xl flex-col gap-4">
        <div className="rounded-lg border border-border bg-card p-4">
          <p className="text-card-title font-medium text-white">{movie.title}</p>
          <p className="mt-1 text-body text-muted">
            {hall.cinema_name || hall.name} · {hall.city || session.city || '—'}
          </p>
          <p className="mt-2 text-body text-muted">
            Seats:{' '}
            {selectedSeats.map((s) => `${s.row}${s.number}`).join(', ') || '—'}
          </p>
          <p className="mt-2 text-body font-medium text-white">
            Total {formatPrice(total)}
          </p>
        </div>

        <TicketCategoryPicker
          basePrice={session.price}
          value={ticketCategory}
          onChange={(id: TicketCategoryId) => setTicketCategory(id)}
        />

        <div className="flex flex-wrap gap-2">
          {tabs.map((t) => (
            <button
              key={t.id}
              type="button"
              onClick={() => setTab(t.id)}
              className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${
                tab === t.id
                  ? 'border-accent bg-accentDim text-accent'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>
        {tab === 'wallet' ? (
          <div className="rounded-lg border border-border bg-card2 p-4 text-body">
            <p className="text-muted">Balance</p>
            <p className="mt-1 text-white">{formatPrice(balance)}</p>
            {sufficient ? (
              <p className="mt-3 text-[#22c55e]">Sufficient balance</p>
            ) : (
              <div className="mt-3 flex flex-wrap items-center gap-2">
                <p className="text-danger">Insufficient balance</p>
                <Button
                  type="button"
                  variant="secondary"
                  className="px-2 py-1 text-[11px]"
                  onClick={() => setTopUpOpen(true)}
                >
                  Top up wallet
                </Button>
              </div>
            )}
          </div>
        ) : (
          <p className="text-body text-muted">
            {tab === 'card' ? 'Card payment (demo).' : 'Kaspi QR (demo).'}
          </p>
        )}
        <Button
          type="button"
          disabled={payDisabled}
          className="w-full max-w-xs gap-2 self-start"
          onClick={() => void handleConfirm()}
        >
          {loading ? <Spinner tone="onPrimary" /> : null}
          Confirm & pay
        </Button>
      </div>
      <TopUpModal
        open={topUpOpen}
        currentBalance={balance}
        onClose={() => setTopUpOpen(false)}
        onConfirm={async (amount) => {
          try {
            const updated = await topUpWallet(amount)
            setUser(updated)
            setTopUpOpen(false)
          } catch {
            onError('Could not top up balance.')
          }
        }}
      />
    </>
  )
}
