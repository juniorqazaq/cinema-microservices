import type { ReactNode } from 'react'

type Tone = 'neutral' | 'success' | 'warning' | 'danger' | 'accent'

const tones: Record<Tone, string> = {
  neutral: 'border-border bg-card2 text-muted',
  success: 'border-[#166534] bg-[#052210] text-[#22c55e]',
  warning: 'border-sky-900/40 bg-accentDim text-vip',
  danger: 'border-danger/40 bg-dangerDim text-danger',
  accent: 'border-accent/50 bg-accentDim text-accent',
}

export function Badge({
  tone = 'neutral',
  children,
  className = '',
}: {
  tone?: Tone
  children: ReactNode
  className?: string
}) {
  return (
    <span
      className={`inline-flex items-center rounded-md border px-2 py-0.5 text-[11px] font-medium ${tones[tone]} ${className}`}
    >
      {children}
    </span>
  )
}
