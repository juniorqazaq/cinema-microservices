import { useEffect, useState } from 'react'
import {
  NavLink,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
import { MapPin, Search } from 'lucide-react'
import { useAuthStore } from '../../store/authStore'
import { useAuth } from '../../hooks/useAuth'
import { Button } from '../ui/Button'
import { MovieHouseLogo } from '../brand/MovieHouseLogo'

function initials(email: string): string {
  const clean = email.replace(/[^a-zA-Z0-9]/g, '')
  return (clean.slice(0, 2) || email.slice(0, 2)).toUpperCase()
}

export function Navbar() {
  const navigate = useNavigate()
  const location = useLocation()
  const [params, setParams] = useSearchParams()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const user = useAuthStore((s) => s.user)
  const { logoutMutation } = useAuth()
  const qParam = params.get('q') ?? ''
  const [searchVal, setSearchVal] = useState(qParam)

  useEffect(() => {
    setSearchVal(qParam)
  }, [qParam])

  async function handleLogout() {
    await logoutMutation.mutateAsync()
    navigate('/movies')
  }

  function applySearch(raw: string) {
    const q = raw.trim()
    if (location.pathname === '/movies') {
      if (q) setParams({ q }, { replace: true })
      else setParams({}, { replace: true })
    } else {
      navigate({
        pathname: '/movies',
        search: q ? `?q=${encodeURIComponent(q)}` : '',
      })
    }
  }

  const navLinkClass = ({ isActive }: { isActive: boolean }) =>
    `whitespace-nowrap rounded-lg px-1 py-1 text-body font-medium transition-colors duration-150 ${
      isActive ? 'text-white' : 'text-muted hover:text-white'
    }`

  return (
    <header className="border-b border-border bg-page">
      <div className="mx-auto max-w-6xl px-4 py-3">
        <div className="flex flex-wrap items-center gap-3 md:flex-nowrap">
          <MovieHouseLogo variant="navbar" className="shrink-0" />

          <div className="ml-auto flex shrink-0 items-center gap-2">
            <Button
              type="button"
              variant="secondary"
              className="hidden rounded-full px-3 py-1.5 text-body sm:inline-flex"
              onClick={() => {}}
            >
              <MapPin className="h-4 w-4 shrink-0" aria-hidden />
              <span className="ml-1">Location</span>
            </Button>
            {!isAuthenticated ? (
              <>
                <Button
                  variant="ghost"
                  type="button"
                  className="rounded-full px-3 py-1.5"
                  onClick={() => navigate('/login')}
                >
                  Sign in
                </Button>
                <Button
                  type="button"
                  className="rounded-full px-3 py-1.5"
                  onClick={() => navigate('/register')}
                >
                  Register
                </Button>
              </>
            ) : (
              <>
                {user?.role === 'admin' ? (
                  <NavLink
                    to="/admin"
                    className="hidden rounded-lg px-2 py-1 text-body font-medium text-muted hover:text-white sm:inline-block"
                  >
                    Admin
                  </NavLink>
                ) : null}
                <div
                  className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-accent bg-accentDim text-body font-medium text-accent"
                  title={user?.email}
                >
                  {user?.email ? initials(user.email) : '—'}
                </div>
                <Button
                  variant="secondary"
                  type="button"
                  className="hidden rounded-full px-3 py-1.5 sm:inline-flex"
                  onClick={() => void handleLogout()}
                  disabled={logoutMutation.isPending}
                >
                  Sign out
                </Button>
              </>
            )}
          </div>
        </div>

        <div className="mt-3 flex flex-col gap-3 border-t border-border pt-3 sm:flex-row sm:items-center sm:gap-4">
          <nav className="flex shrink-0 flex-wrap items-center gap-x-6 gap-y-1 text-body">
            <NavLink to="/movies" className={navLinkClass} end>
              Home
            </NavLink>
            <NavLink to="/profile" className={navLinkClass}>
              My tickets
            </NavLink>
            {isAuthenticated && user?.role === 'admin' ? (
              <NavLink
                to="/admin"
                className={({ isActive }) =>
                  `${navLinkClass({ isActive })} sm:hidden`
                }
              >
                Admin
              </NavLink>
            ) : null}
            {isAuthenticated ? (
              <Button
                variant="ghost"
                type="button"
                className="rounded-full px-2 py-1 sm:hidden"
                onClick={() => void handleLogout()}
                disabled={logoutMutation.isPending}
              >
                Sign out
              </Button>
            ) : null}
          </nav>

          <form
            className="min-w-0 flex-1"
            onSubmit={(e) => {
              e.preventDefault()
              applySearch(searchVal)
            }}
          >
            <div className="relative w-full">
              <Search
                className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted"
                aria-hidden
              />
              <input
                type="search"
                name="q"
                placeholder="Search movies…"
                value={searchVal}
                onChange={(e) => setSearchVal(e.target.value)}
                className="w-full rounded-full border border-border2 bg-card py-2 pl-9 pr-4 text-body text-white outline-none transition-colors placeholder:text-muted focus:border-accent"
                aria-label="Search movies"
              />
            </div>
          </form>
        </div>
      </div>
    </header>
  )
}
