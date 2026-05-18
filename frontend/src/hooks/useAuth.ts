import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  changePassword,
  loginUser,
  logoutUser,
  registerUser,
} from '../api/auth'
import { useAuthStore } from '../store/authStore'
import { fetchProfile } from '../api/wallet'
import { userFromAccessToken } from '../utils/jwt'

export function useAuth() {
  const queryClient = useQueryClient()
  const user = useAuthStore((s) => s.user)
  const setTokens = useAuthStore((s) => s.setTokens)
  const setUser = useAuthStore((s) => s.setUser)
  const logoutStore = useAuthStore((s) => s.logout)

  const loginMutation = useMutation({
    mutationFn: async (vars: { email: string; password: string }) => {
      const tokens = await loginUser(vars.email, vars.password)
      setTokens(tokens.access_token, tokens.refresh_token)
      try {
        setUser(await fetchProfile())
      } catch {
        setUser(userFromAccessToken(tokens.access_token, vars.email))
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bookings'] })
    },
  })

  const registerMutation = useMutation({
    mutationFn: async (vars: { email: string; password: string }) => {
      const tokens = await registerUser(vars.email, vars.password)
      setTokens(tokens.access_token, tokens.refresh_token)
      try {
        setUser(await fetchProfile())
      } catch {
        setUser(userFromAccessToken(tokens.access_token, vars.email))
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bookings'] })
    },
  })

  const logoutMutation = useMutation({
    mutationFn: async () => {
      try {
        await logoutUser()
      } finally {
        logoutStore()
        queryClient.clear()
      }
    },
  })

  const changePasswordMutation = useMutation({
    mutationFn: (vars: { old_password: string; new_password: string }) =>
      changePassword(vars.old_password, vars.new_password),
  })

  return {
    user: user ?? null,
    loginMutation,
    registerMutation,
    logoutMutation,
    changePasswordMutation,
  }
}
