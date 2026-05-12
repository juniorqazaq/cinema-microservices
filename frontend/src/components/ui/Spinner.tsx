import { Loader2 } from 'lucide-react'

type SpinnerTone = 'accent' | 'onPrimary'

export function Spinner({
  className = '',
  tone = 'accent',
}: {
  className?: string
  tone?: SpinnerTone
}) {
  const color = tone === 'onPrimary' ? 'text-white' : 'text-accent'
  return (
    <Loader2
      className={`h-4 w-4 shrink-0 animate-spin ${color} ${className}`}
      aria-hidden="true"
    />
  )
}
