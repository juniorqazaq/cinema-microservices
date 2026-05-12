import { useAdminStats } from '../../hooks/useBooking'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import { getErrorMessage } from '../../utils/errorHandler'

export function AdminDashboardPage() {
  const statsQuery = useAdminStats()

  if (statsQuery.isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner className="h-6 w-6" />
      </div>
    )
  }

  if (statsQuery.isError || !statsQuery.data) {
    return <ErrorBanner message={getErrorMessage(statsQuery.error)} />
  }

  const s = statsQuery.data

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-section font-light text-white">Dashboard</h1>
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-lg border border-border bg-card p-4">
          <p className="text-body text-muted">Total</p>
          <p className="mt-1 text-hero font-light text-white">{s.total}</p>
        </div>
        <div className="rounded-lg border border-border bg-card p-4">
          <p className="text-body text-muted">Confirmed</p>
          <p className="mt-1 text-hero font-light text-white">{s.confirmed}</p>
        </div>
        <div className="rounded-lg border border-border bg-card p-4">
          <p className="text-body text-muted">Cancelled</p>
          <p className="mt-1 text-hero font-light text-white">{s.cancelled}</p>
        </div>
      </div>
    </div>
  )
}
