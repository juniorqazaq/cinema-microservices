import { get, post } from './axios'
import type { User } from '../types'

export async function fetchProfile(): Promise<User> {
  return get<User>('/users/me')
}

export async function topUpWallet(amount: number): Promise<User> {
  return post<User>('/wallet/topup', { amount })
}

export async function getCities(): Promise<string[]> {
  return get<string[]>('/cities')
}
