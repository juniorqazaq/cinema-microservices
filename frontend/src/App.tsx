import { Navigate, Route, Routes } from 'react-router-dom'
import { AppShell } from './components/layout/AppShell'
import { AdminLayout } from './components/admin/AdminLayout'
import { ProtectedRoute } from './components/layout/ProtectedRoute'
import { AdminRoute } from './components/layout/AdminRoute'
import { LoginPage } from './pages/auth/LoginPage'
import { RegisterPage } from './pages/auth/RegisterPage'
import { HomePage } from './pages/movies/HomePage'
import { MovieDetailPage } from './pages/movies/MovieDetailPage'
import { BookingPage } from './pages/booking/BookingPage'
import { BookingSuccessPage } from './pages/booking/BookingSuccessPage'
import { ProfilePage } from './pages/profile/ProfilePage'
import { AdminDashboardPage } from './pages/admin/AdminDashboardPage'
import { AdminMoviesPage } from './pages/admin/AdminMoviesPage'
import { AdminUsersPage } from './pages/admin/AdminUsersPage'
import { AdminBookingsPage } from './pages/admin/AdminBookingsPage'

export default function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route index element={<Navigate to="/movies" replace />} />
        <Route path="movies" element={<HomePage />} />
        <Route path="movies/:id" element={<MovieDetailPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="booking" element={<BookingPage />} />
          <Route path="booking/success" element={<BookingSuccessPage />} />
          <Route path="profile" element={<ProfilePage />} />
          <Route element={<AdminRoute />}>
            <Route path="admin" element={<AdminLayout />}>
              <Route index element={<AdminDashboardPage />} />
              <Route path="movies" element={<AdminMoviesPage />} />
              <Route path="users" element={<AdminUsersPage />} />
              <Route path="bookings" element={<AdminBookingsPage />} />
            </Route>
          </Route>
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/movies" replace />} />
    </Routes>
  )
}
