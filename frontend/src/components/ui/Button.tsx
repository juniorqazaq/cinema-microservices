import type { ButtonHTMLAttributes, ReactNode } from 'react'

type Variant = 'primary' | 'secondary' | 'ghost'

const variants: Record<Variant, string> = {
  primary:
    'bg-accent !text-white border border-accent hover:bg-accentHover hover:border-accentHover',
  secondary:
    'bg-card2 !text-white border border-border hover:border-accent',
  ghost:
    'bg-transparent !text-white border border-transparent hover:border-border2',
}

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  children: ReactNode
}

export function Button({
  variant = 'primary',
  className = '',
  disabled,
  children,
  type = 'button',
  ...rest
}: ButtonProps) {
  return (
    <button
      {...rest}
      type={type}
      disabled={disabled}
      className={`inline-flex items-center justify-center gap-2 rounded-lg border px-4 py-2 text-body font-medium transition-colors duration-150 ease-out disabled:cursor-not-allowed disabled:opacity-50 ${variants[variant]} ${className}`}
    >
      {children}
    </button>
  )
}
