import { X } from 'lucide-react'

interface ErrorBannerProps {
  message: string
  onDismiss?: () => void
}

export function ErrorBanner({ message, onDismiss }: ErrorBannerProps) {
  return (
    <div
      className="flex items-start justify-between gap-3 rounded-lg border border-danger/50 bg-dangerDim px-4 py-3 text-body text-danger"
      role="alert"
    >
      <p className="min-w-0 flex-1">{message}</p>
      {onDismiss ? (
        <button
          type="button"
          className="shrink-0 rounded border border-transparent p-0.5 text-danger hover:border-border2"
          aria-label="Dismiss"
          onClick={onDismiss}
        >
          <X className="h-4 w-4" aria-hidden="true" />
        </button>
      ) : null}
    </div>
  )
}
