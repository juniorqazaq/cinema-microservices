import type { ReactNode } from 'react'

export function SectionTitle({
  children,
  action,
}: {
  children: ReactNode
  action?: ReactNode
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <div className="flex items-center gap-3">
        <span
          className="h-6 w-1 shrink-0 rounded-sm bg-accent"
          aria-hidden
        />
        <h2 className="text-section font-medium text-white">{children}</h2>
      </div>
      {action}
    </div>
  )
}
