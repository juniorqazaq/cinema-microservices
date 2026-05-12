import { Outlet, useLocation } from 'react-router-dom'
import { Navbar } from './Navbar'

const pathsWithoutNavbar = ['/login', '/register']

export function AppShell() {
  const { pathname } = useLocation()
  const hideNavbar = pathsWithoutNavbar.includes(pathname)

  return (
    <div className="flex min-h-screen flex-col bg-page">
      {!hideNavbar ? <Navbar /> : null}
      <main
        className={
          hideNavbar
            ? 'mx-auto flex w-full max-w-6xl flex-1 flex-col items-center justify-center px-4 py-6'
            : 'mx-auto w-full max-w-6xl flex-1 px-4 py-6'
        }
      >
        <Outlet />
      </main>
    </div>
  )
}
