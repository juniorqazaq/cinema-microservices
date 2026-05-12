import {
  forwardRef,
  type InputHTMLAttributes,
  type ReactNode,
  type SelectHTMLAttributes,
} from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  id: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, id, className = '', ...rest }, ref) => {
    return (
      <div className="flex flex-col gap-1 text-left">
        {label ? (
          <label htmlFor={id} className="text-body font-medium text-white">
            {label}
          </label>
        ) : null}
        <input
          ref={ref}
          id={id}
          className={`rounded-lg border border-border2 bg-card2 px-3 py-2 text-body text-white outline-none transition-colors duration-150 placeholder:text-muted focus:border-accent ${className}`}
          {...rest}
        />
        {error ? (
          <p className="text-body text-danger" role="alert">
            {error}
          </p>
        ) : null}
      </div>
    )
  },
)

Input.displayName = 'Input'

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string
  error?: string
  id: string
  children: ReactNode
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, error, id, className = '', children, ...rest }, ref) => {
    return (
      <div className="flex flex-col gap-1 text-left">
        {label ? (
          <label htmlFor={id} className="text-body font-medium text-white">
            {label}
          </label>
        ) : null}
        <select
          ref={ref}
          id={id}
          className={`rounded-lg border border-border2 bg-card2 px-3 py-2 text-body text-white outline-none transition-colors duration-150 focus:border-accent ${className}`}
          {...rest}
        >
          {children}
        </select>
        {error ? (
          <p className="text-body text-danger" role="alert">
            {error}
          </p>
        ) : null}
      </div>
    )
  },
)

Select.displayName = 'Select'
