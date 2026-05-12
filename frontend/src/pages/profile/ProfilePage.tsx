import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { useQueries } from '@tanstack/react-query'
import { isFuture, parseISO } from 'date-fns'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Ticket } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'
import { useBookingHistory } from '../../hooks/useBooking'
import { getSession } from '../../api/movies'
import { TicketCard } from '../../components/profile/TicketCard'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Input } from '../../components/ui/Input'
import { Spinner } from '../../components/ui/Spinner'
import { useAuthStore } from '../../store/authStore'
import { formatDate } from '../../utils/format'
import { getErrorMessage } from '../../utils/errorHandler'
import type { Booking, Session } from '../../types'

const passwordSchema = z
  .object({
    old_password: z.string().min(1, 'Required'),
    new_password: z.string().min(8, 'Use at least 8 characters'),
    confirm: z.string().min(1, 'Required'),
  })
  .refine((d) => d.new_password === d.confirm, {
    message: 'Passwords do not match',
    path: ['confirm'],
  })

type PasswordForm = z.infer<typeof passwordSchema>

type Tab = 'upcoming' | 'past' | 'all'

function bookingCategory(
  booking: Booking,
  sessionStart: string | undefined,
): 'upcoming' | 'past' {
  if (booking.status === 'cancelled') return 'past'
  if (!sessionStart) return 'past'
  const d = parseISO(sessionStart)
  if (booking.status === 'confirmed' && isFuture(d)) return 'upcoming'
  return 'past'
}

export function ProfilePage() {
  const user = useAuthStore((s) => s.user)
  const { changePasswordMutation } = useAuth()
  const historyQuery = useBookingHistory()
  const [tab, setTab] = useState<Tab>('upcoming')
  const [pwOpen, setPwOpen] = useState(false)

  const bookings = historyQuery.data ?? []

  const sessionQueries = useQueries({
    queries: bookings.map((b) => ({
      queryKey: ['session', b.session_id],
      queryFn: () => getSession(b.session_id),
      enabled: Boolean(b.session_id) && historyQuery.isSuccess,
    })),
  })

  const sessionByBookingId = useMemo(() => {
    const map: Record<string, Session | undefined> = {}
    bookings.forEach((b, i) => {
      map[b.id] = sessionQueries[i]?.data
    })
    return map
  }, [bookings, sessionQueries])

  const filtered = useMemo(() => {
    if (tab === 'all') return bookings
    return bookings.filter((b) => {
      const cat = bookingCategory(b, sessionByBookingId[b.id]?.start_time)
      return tab === 'upcoming' ? cat === 'upcoming' : cat === 'past'
    })
  }, [bookings, tab, sessionByBookingId])

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      old_password: '',
      new_password: '',
      confirm: '',
    },
  })

  if (!user) {
    return (
      <div className="rounded-lg border border-border bg-card p-4 text-body text-muted">
        Not found
      </div>
    )
  }

  const tabs: { id: Tab; label: string }[] = [
    { id: 'upcoming', label: 'Upcoming' },
    { id: 'past', label: 'Past' },
    { id: 'all', label: 'All' },
  ]

  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-8">
      <section className="rounded-lg border border-border bg-card p-4">
        <div className="flex flex-wrap items-start gap-4">
          <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full border border-accent bg-accentDim text-body font-medium text-accent">
            {(user.email.slice(0, 2) || '??').toUpperCase()}
          </div>
          <div className="min-w-0">
            <p className="text-card-title font-medium text-white">{user.email}</p>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              {user.role === 'admin' ? (
                <Badge tone="accent">admin</Badge>
              ) : (
                <Badge tone="neutral">user</Badge>
              )}
              <span className="text-body text-muted">
                Member since:{' '}
                {user.created_at ? formatDate(user.created_at) : '—'}
              </span>
            </div>
          </div>
        </div>
      </section>

      <section className="rounded-lg border border-border bg-card p-4">
        <div className="flex flex-wrap gap-2 border-b border-border pb-3">
          {tabs.map((t) => (
            <button
              key={t.id}
              type="button"
              onClick={() => setTab(t.id)}
              className={`rounded-lg border px-3 py-1.5 text-body font-medium transition-colors duration-150 ${
                tab === t.id
                  ? 'border-accent bg-accentDim text-accent'
                  : 'border-border bg-card2 text-muted hover:border-border2'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>
        <div className="mt-4 flex flex-col gap-3">
          {historyQuery.isLoading ? (
            <div className="flex justify-center py-8">
              <Spinner />
            </div>
          ) : historyQuery.isError ? (
            <ErrorBanner message={getErrorMessage(historyQuery.error)} />
          ) : bookings.length === 0 ? (
            <div className="flex flex-col items-center gap-3 py-10 text-center">
              <Ticket className="h-10 w-10 text-muted" aria-hidden />
              <p className="text-body text-muted">No tickets yet</p>
              <Link
                to="/movies"
                className="inline-flex items-center justify-center rounded-lg border border-accent bg-accent px-4 py-2 text-body font-medium text-white hover:border-accentHover hover:bg-accentHover"
              >
                Browse movies
              </Link>
            </div>
          ) : filtered.length === 0 ? (
            <p className="py-6 text-center text-body text-muted">
              Nothing in this tab.
            </p>
          ) : (
            filtered.map((b) => (
              <TicketCard
                key={b.id}
                booking={b}
                session={sessionByBookingId[b.id]}
              />
            ))
          )}
        </div>
      </section>

      <section className="rounded-lg border border-border bg-card">
        <button
          type="button"
          className="flex w-full items-center justify-between px-4 py-3 text-left text-body font-medium text-white"
          onClick={() => setPwOpen((v) => !v)}
          aria-expanded={pwOpen}
        >
          Change password
          <span className="text-muted">{pwOpen ? '−' : '+'}</span>
        </button>
        {pwOpen ? (
          <div className="border-t border-border p-4">
            {changePasswordMutation.isError ? (
              <ErrorBanner
                message={getErrorMessage(changePasswordMutation.error)}
              />
            ) : null}
            {changePasswordMutation.isSuccess ? (
              <p className="mb-3 text-body text-[#22c55e]">Password updated.</p>
            ) : null}
            <form
              className="mt-2 flex flex-col gap-3"
              onSubmit={handleSubmit(async (values) => {
                await changePasswordMutation.mutateAsync({
                  old_password: values.old_password,
                  new_password: values.new_password,
                })
                reset()
              })}
            >
              <Input
                id="old-password"
                label="Old password"
                type="password"
                autoComplete="current-password"
                {...register('old_password')}
                error={errors.old_password?.message}
              />
              <Input
                id="new-password"
                label="New password"
                type="password"
                autoComplete="new-password"
                {...register('new_password')}
                error={errors.new_password?.message}
              />
              <Input
                id="confirm-new"
                label="Confirm new password"
                type="password"
                autoComplete="new-password"
                {...register('confirm')}
                error={errors.confirm?.message}
              />
              <Button
                type="submit"
                disabled={changePasswordMutation.isPending}
                className="w-fit gap-2"
              >
                {changePasswordMutation.isPending ? (
                  <Spinner tone="onPrimary" />
                ) : null}
                Update password
              </Button>
            </form>
          </div>
        ) : null}
      </section>
    </div>
  )
}
