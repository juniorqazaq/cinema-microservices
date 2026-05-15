const STORAGE_PREFIX = 'cinema-wallet:'

export const DEFAULT_WALLET_BALANCE = 10_000

export function readWalletBalance(userId: string): number {
  if (!userId) return DEFAULT_WALLET_BALANCE
  try {
    const raw = localStorage.getItem(`${STORAGE_PREFIX}${userId}`)
    if (raw == null) return DEFAULT_WALLET_BALANCE
    const n = Number.parseFloat(raw)
    return Number.isFinite(n) && n >= 0 ? n : DEFAULT_WALLET_BALANCE
  } catch {
    return DEFAULT_WALLET_BALANCE
  }
}

export function writeWalletBalance(userId: string, balance: number): void {
  if (!userId) return
  try {
    localStorage.setItem(`${STORAGE_PREFIX}${userId}`, String(Math.max(0, balance)))
  } catch {
    /* ignore quota errors */
  }
}
