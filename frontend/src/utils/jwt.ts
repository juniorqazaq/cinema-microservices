import type { User } from '../types'

function decodeBase64Url(payload: string): string {
  const pad = payload.length % 4
  const b64 =
    payload.replace(/-/g, '+').replace(/_/g, '/') +
    (pad ? '='.repeat(4 - pad) : '')
  return atob(b64)
}

export function decodeJwtPayload(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    const payload = parts[1]
    if (!payload) return null
    return JSON.parse(decodeBase64Url(payload)) as Record<string, unknown>
  } catch {
    return null
  }
}

export function userFromAccessToken(
  accessToken: string,
  emailFallback?: string,
): User {
  const p = decodeJwtPayload(accessToken) ?? {}
  const raw = String(p.role ?? '')
  const role =
    raw === 'ADMIN' || raw.toLowerCase() === 'admin' ? 'admin' : 'user'
  const id = String(p.sub ?? p.user_id ?? p.id ?? '')
  const email = String(p.email ?? emailFallback ?? '')
  const bal = p.balance
  return {
    id,
    email,
    role,
    is_banned: Boolean(p.is_banned),
    created_at: String(p.created_at ?? ''),
    updated_at: String(p.updated_at ?? ''),
    balance: typeof bal === 'number' ? bal : undefined,
  }
}
