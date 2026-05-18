import { Button } from '../ui/Button'

interface AgeVerifyModalProps {
  open: boolean
  onClose: () => void
  onConfirm: () => void
}

export function AgeVerifyModal({ open, onClose, onConfirm }: AgeVerifyModalProps) {
  if (!open) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="age-verify-title"
    >
      <div className="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-xl">
        <h2 id="age-verify-title" className="text-card-title font-medium text-white">
          18+ confirmation
        </h2>
        <p className="mt-2 text-body text-muted">
          This film is rated 18+. Confirm that you are at least 18 years old to continue.
        </p>
        <div className="mt-6 flex flex-wrap justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="button" onClick={onConfirm}>
            I am 18+
          </Button>
        </div>
      </div>
    </div>
  )
}
