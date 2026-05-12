import { useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import axios from 'axios'
import { ArrowLeft, Eye, EyeOff } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'
import { useAuthStore } from '../../store/authStore'
import { MovieHouseLogo } from '../../components/brand/MovieHouseLogo'
import { Button } from '../../components/ui/Button'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Input } from '../../components/ui/Input'
import { Spinner } from '../../components/ui/Spinner'

const schema = z
  .object({
    email: z.string().email(),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

type FormValues = z.infer<typeof schema>

export function RegisterPage() {
  const navigate = useNavigate()
  const authed = useAuthStore((s) => s.isAuthenticated)
  const { registerMutation } = useAuth()
  const [showPassword, setShowPassword] = useState(false)
  const [banner, setBanner] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: '', password: '', confirmPassword: '' },
  })

  if (authed) {
    return <Navigate to="/movies" replace />
  }

  return (
    <div className="mx-auto w-full max-w-md">
      <Link
        to="/movies"
        className="mb-4 inline-flex items-center gap-2 text-body text-muted transition-colors hover:text-white"
      >
        <ArrowLeft className="h-4 w-4 shrink-0" aria-hidden="true" />
        Back to home
      </Link>
      <div className="rounded-lg border border-border bg-card p-6">
        <div className="mb-6 flex justify-center">
          <MovieHouseLogo variant="auth" to="/movies" />
        </div>
        {banner ? (
          <div className="mb-4">
            <ErrorBanner message={banner} onDismiss={() => setBanner(null)} />
          </div>
        ) : null}
        <form
          className="flex flex-col gap-4"
          onSubmit={handleSubmit(async (values) => {
            setBanner(null)
            try {
              await registerMutation.mutateAsync({
                email: values.email,
                password: values.password,
              })
              navigate('/movies', { replace: true })
            } catch (e) {
              if (axios.isAxiosError(e) && e.response?.status === 409) {
                setBanner('An account with this email already exists')
              }
            }
          })}
        >
          <Input
            id="register-email"
            label="Email"
            type="email"
            autoComplete="email"
            {...register('email')}
            error={errors.email?.message}
          />
          <div className="relative">
            <Input
              id="register-password"
              label="Password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="new-password"
              {...register('password')}
              error={errors.password?.message}
            />
            <button
              type="button"
              className="absolute right-2 top-9 rounded border border-transparent p-1 text-muted hover:border-border2"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
              onClick={() => setShowPassword((v) => !v)}
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4" aria-hidden="true" />
              ) : (
                <Eye className="h-4 w-4" aria-hidden="true" />
              )}
            </button>
          </div>
          <Input
            id="register-confirm"
            label="Confirm password"
            type={showPassword ? 'text' : 'password'}
            autoComplete="new-password"
            {...register('confirmPassword')}
            error={errors.confirmPassword?.message}
          />
          <Button
            type="submit"
            disabled={registerMutation.isPending}
            className="w-full gap-2 rounded-full border-white/25 py-3 text-[14px] font-semibold shadow-sm"
          >
            {registerMutation.isPending ? <Spinner tone="onPrimary" /> : null}
            Register
          </Button>
        </form>
        <p className="mt-4 text-center text-body text-muted">
          <Link to="/login" className="text-accent hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
