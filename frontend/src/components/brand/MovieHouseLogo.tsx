import { Link } from 'react-router-dom'

/** White film-strip mark (matches MOVIE HOUSE reference). Uses `currentColor`. */
export function MovieHouseMark({ className = 'h-9 w-9' }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 40 44"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden
    >
      <rect
        x="1"
        y="1"
        width="38"
        height="42"
        rx="4"
        stroke="currentColor"
        strokeWidth="2"
      />
      {[10, 18, 26].map((cy) => (
        <rect key={`l${cy}`} x="5" y={cy} width="4" height="4" rx="0.5" fill="currentColor" />
      ))}
      {[10, 18, 26].map((cy) => (
        <rect key={`r${cy}`} x="31" y={cy} width="4" height="4" rx="0.5" fill="currentColor" />
      ))}
      <rect x="13" y="14" width="14" height="12" rx="1" stroke="currentColor" strokeWidth="1.5" />
    </svg>
  )
}

type LogoVariant = 'navbar' | 'auth' | 'footer'

const textByVariant: Record<LogoVariant, string> = {
  navbar: 'text-card-title font-medium tracking-wide',
  auth: 'text-lg font-medium md:text-xl',
  footer: 'text-body font-medium',
}

const stackByVariant: Record<LogoVariant, string> = {
  navbar: 'flex flex-row items-center gap-2.5 text-white',
  auth: 'flex flex-col items-center gap-2 text-white',
  footer: 'flex flex-row items-center gap-2 text-white',
}

const markByVariant: Record<LogoVariant, string> = {
  navbar: 'h-8 w-[29px] shrink-0',
  auth: 'h-11 w-[40px] shrink-0',
  footer: 'h-7 w-[26px] shrink-0 text-white/90',
}

interface MovieHouseLogoProps {
  variant?: LogoVariant
  to?: string
  className?: string
}

export function MovieHouseLogo({
  variant = 'navbar',
  to = '/movies',
  className = '',
}: MovieHouseLogoProps) {
  const inner = (
    <>
      <MovieHouseMark className={`${markByVariant[variant]} text-white`} />
      <span className={`${textByVariant[variant]} text-white`}>Movie house</span>
    </>
  )

  return (
    <Link
      to={to}
      className={`${stackByVariant[variant]} transition-opacity hover:opacity-90 ${className}`}
    >
      {inner}
    </Link>
  )
}
