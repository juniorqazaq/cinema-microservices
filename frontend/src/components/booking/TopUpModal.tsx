import { useState } from 'react'
import { formatPrice } from '../../utils/format'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'

const PRESETS = [5_000, 10_000, 20_000] as const

interface TopUpModalProps {
  open: boolean
  currentBalance: number
  onClose: () => void
  onConfirm: (amount: number) => void
}

export function TopUpModal({
  open,
  currentBalance,
  onClose,
  onConfirm,
}: TopUpModalProps) {
  const [custom, setCustom] = useState('')

  if (!open) return null

  function apply(amount: number) {
    if (amount > 0) onConfirm(amount)
    setCustom('')
    onClose()
  }

  const customAmount = Number.parseFloat(custom)
  const customValid = Number.isFinite(customAmount) && customAmount >= 100

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        className="absolute inset-0 bg-black/70"
        aria-label="Close"
        onClick={onClose}
      />
      <div
        className="relative z-10 w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-xl"
        role="dialog"
        aria-modal="true"
        aria-labelledby="topup-title"
      >
        <h2 id="topup-title" className="text-section font-medium text-white">
          Top up wallet
        </h2>
        <p className="mt-3 text-body text-muted">
          Current balance:{' '}
          <span className="font-medium text-white">{formatPrice(currentBalance)}</span>
        </p>
        <div className="mt-4 grid grid-cols-3 gap-2">
          {PRESETS.map((amount) => (
            <button
              key={amount}
              type="button"
              onClick={() => apply(amount)}
              className="rounded-lg border border-border bg-card2 px-3 py-2 text-body font-medium text-white transition-colors hover:border-accent"
            >
              +{formatPrice(amount)}
            </button>
          ))}
        </div>
        <div className="mt-4">
          <Input
            id="topup-amount"
            label="Custom amount (₸)"
            type="number"
            min={100}
            step={100}
            value={custom}
            onChange={(e) => setCustom(e.target.value)}
          />
        </div>
        <div className="mt-4 flex flex-wrap justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button
            type="button"
            onClick={() => apply(customAmount)}
            disabled={!customValid}
          >
            Add funds
          </Button>
        </div>
      </div>
    </div>
  )
}
