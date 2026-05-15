import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'
import { readWalletBalance, writeWalletBalance } from '../lib/wallet'

export interface AuthState {
  user: User | null
  accessToken: string
  refreshToken: string
  isAuthenticated: boolean
  setTokens: (accessToken: string, refreshToken: string) => void
  setUser: (user: User | null) => void
  topUpBalance: (amount: number) => void
  adjustBalance: (delta: number) => void
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
    (set, get) => ({
      ...empty(),
      setTokens: (accessToken, refreshToken) =>
        set({
          accessToken,
          refreshToken,
          isAuthenticated: Boolean(accessToken && refreshToken),
        }),
      setUser: (user) => {
        if (!user) {
          set({ user: null })
          return
        }
        const balance = user.balance ?? readWalletBalance(user.id)
        writeWalletBalance(user.id, balance)
        set({ user: { ...user, balance } })
      },
      topUpBalance: (amount) => {
        const { user } = get()
        if (!user || amount <= 0) return
        const balance = (user.balance ?? 0) + amount
        writeWalletBalance(user.id, balance)
        set({ user: { ...user, balance } })
      },
      adjustBalance: (delta) => {
        const { user } = get()
        if (!user) return
        const balance = Math.max(0, (user.balance ?? 0) + delta)
        writeWalletBalance(user.id, balance)
        set({ user: { ...user, balance } })
      },
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
