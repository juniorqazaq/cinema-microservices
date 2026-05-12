import axios, {
  type AxiosError,
  type AxiosInstance,
  type InternalAxiosRequestConfig,
} from 'axios'
import { useAuthStore } from '../store/authStore'
import type { AuthTokens, RefreshResponse } from '../types'

export const BASE_URL =
  import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

type ApiEnvelope<T> = { data: T }
type ApiErrorBody = { error: string }

const rawClient = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

let isRefreshing = false
let queue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

function flushQueue(error: unknown, token: string | null) {
  queue.forEach((p) => {
    if (error) {
      p.reject(error)
    } else if (token) {
      p.resolve(token)
    }
  })
  queue = []
}

async function refreshAccessToken(): Promise<string> {
  const refreshToken = useAuthStore.getState().refreshToken
  if (!refreshToken) {
    throw new Error('missing refresh token')
  }
  const res = await rawClient.post<ApiEnvelope<RefreshResponse>>(
    '/auth/refresh',
    { refresh_token: refreshToken },
  )
  const access_token = res.data.data.access_token
  const currentRefresh = useAuthStore.getState().refreshToken
  useAuthStore.getState().setTokens(access_token, currentRefresh)
  return access_token
}

export const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().accessToken
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiErrorBody>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
    }
    if (!originalRequest) {
      return Promise.reject(error)
    }

    const status = error.response?.status
    const url = String(originalRequest.url ?? '')

    if (status === 401) {
      if (
        url.includes('/auth/login') ||
        url.includes('/auth/register') ||
        url.includes('/auth/refresh')
      ) {
        if (url.includes('/auth/refresh')) {
          useAuthStore.getState().logout()
          window.location.assign('/login')
        }
        return Promise.reject(error)
      }

      if (originalRequest._retry) {
        useAuthStore.getState().logout()
        window.location.assign('/login')
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          queue.push({
            resolve: (token: string) => {
              originalRequest.headers.Authorization = `Bearer ${token}`
              originalRequest._retry = true
              resolve(api(originalRequest))
            },
            reject,
          })
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const newAccess = await refreshAccessToken()
        flushQueue(null, newAccess)
        originalRequest.headers.Authorization = `Bearer ${newAccess}`
        return api(originalRequest)
      } catch (e) {
        flushQueue(e, null)
        useAuthStore.getState().logout()
        window.location.assign('/login')
        return Promise.reject(e)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  },
)

export function unwrapData<T>(payload: ApiEnvelope<T>): T {
  return payload.data
}

export async function get<T>(url: string, config?: object): Promise<T> {
  const res = await api.get<ApiEnvelope<T>>(url, config)
  return unwrapData(res.data)
}

export async function post<T>(url: string, body?: unknown): Promise<T> {
  const res = await api.post(url, body)
  const payload = res.data as { data?: T } | undefined
  if (payload && 'data' in payload && payload.data !== undefined) {
    return payload.data
  }
  return {} as T
}

export async function put<T>(url: string, body?: unknown): Promise<T> {
  const res = await api.put(url, body)
  const payload = res.data as { data?: T } | undefined
  if (payload && 'data' in payload && payload.data !== undefined) {
    return payload.data
  }
  return {} as T
}

export async function del<T>(url: string): Promise<T> {
  const res = await api.delete(url)
  const payload = res.data as { data?: T } | undefined
  if (payload && 'data' in payload && payload.data !== undefined) {
    return payload.data
  }
  return {} as T
}

/** Auth endpoints use rawClient (no Bearer on register/login). */
export async function postAuth<T>(
  path: string,
  body: Record<string, unknown>,
): Promise<T> {
  const res = await rawClient.post<ApiEnvelope<T>>(path, body)
  return unwrapData(res.data)
}

export async function postTokens(
  path: string,
  body: Record<string, unknown>,
): Promise<AuthTokens> {
  return postAuth<AuthTokens>(path, body)
}
