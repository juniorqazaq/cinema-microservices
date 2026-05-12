import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'

export interface AuthState {
  user: User | null
  accessToken: string
  refreshToken: string
  isAuthenticated: boolean
  setTokens: (accessToken: string, refreshToken: string) => void
  setUser: (user: User | null) => void
  logout: () => void
}

const empty = () => ({
  user: null as User | null,
  accessToken: '',
  refreshToken: '',
  isAuthenticated: false,
})

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      ...empty(),
      setTokens: (accessToken, refreshToken) =>
        set({
          accessToken,
          refreshToken,
          isAuthenticated: Boolean(accessToken && refreshToken),
        }),
      setUser: (user) => set({ user }),
      logout: () => set(empty()),
    }),
    {
      name: 'cinema-auth',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
)
