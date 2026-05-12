import { NavLink, Outlet } from 'react-router-dom'
import { Calendar, Film, LayoutDashboard, Users } from 'lucide-react'
import { MovieHouseLogo } from '../brand/MovieHouseLogo'

const linkClass = ({ isActive }: { isActive: boolean }) =>
  `flex items-center gap-2 rounded-lg border px-3 py-2 text-body font-medium transition-colors duration-150 ${
    isActive
      ? 'border-accent bg-accentDim text-accent'
      : 'border-transparent bg-card2 text-muted hover:border-border2'
  }`

export function AdminLayout() {
  return (
    <div className="flex min-h-[calc(100vh-5rem)] flex-col gap-6 lg:flex-row">
      <aside className="w-full shrink-0 rounded-lg border border-border bg-card p-3 lg:w-52">
        <div className="mb-3 border-b border-border pb-3">
          <MovieHouseLogo variant="footer" className="justify-center lg:justify-start" />
        </div>
        <p className="mb-2 px-2 text-[11px] font-medium text-muted">Admin</p>
        <nav className="flex flex-col gap-1">
          <NavLink to="/admin" end className={linkClass}>
            <LayoutDashboard className="h-4 w-4" aria-hidden="true" />
            Dashboard
          </NavLink>
          <NavLink to="/admin/movies" className={linkClass}>
            <Film className="h-4 w-4" aria-hidden="true" />
            Movies
          </NavLink>
          <NavLink to="/admin/users" className={linkClass}>
            <Users className="h-4 w-4" aria-hidden="true" />
            Users
          </NavLink>
          <NavLink to="/admin/bookings" className={linkClass}>
            <Calendar className="h-4 w-4" aria-hidden="true" />
            Bookings
          </NavLink>
        </nav>
      </aside>
      <div className="min-w-0 flex-1">
        <Outlet />
      </div>
    </div>
  )
}
