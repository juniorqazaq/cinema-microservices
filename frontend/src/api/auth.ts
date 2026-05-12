import { post, postAuth, postTokens, put } from './axios'
import type { AuthTokens, RefreshResponse } from '../types'

export async function registerUser(
  email: string,
  password: string,
): Promise<AuthTokens> {
  return postTokens('/auth/register', { email, password })
}

export async function loginUser(
  email: string,
  password: string,
): Promise<AuthTokens> {
  return postTokens('/auth/login', { email, password })
}

export async function logoutUser(): Promise<Record<string, never>> {
  return post<Record<string, never>>('/auth/logout', {})
}

export async function refreshToken(
  refresh_token: string,
): Promise<RefreshResponse> {
  return postAuth<RefreshResponse>('/auth/refresh', { refresh_token })
}

export async function changePassword(
  old_password: string,
  new_password: string,
): Promise<Record<string, never>> {
  return put<Record<string, never>>('/users/password', {
    old_password,
    new_password,
  })
}
