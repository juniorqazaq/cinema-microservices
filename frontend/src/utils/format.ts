import { format, formatDistanceToNow, isFuture, parseISO } from 'date-fns'

export function formatPrice(n: number): string {
  return `₸${n.toLocaleString('ru-RU')}`
}

export function formatDuration(minutes: number): string {
  return `${Math.floor(minutes / 60)}h ${minutes % 60}m`
}

export function formatDate(iso: string): string {
  if (!iso) return '—'
  return format(parseISO(iso), 'EEE, d MMM')
}

const ALMATY_TZ = 'Asia/Almaty'

export function formatTime(iso: string): string {
  if (!iso) return '—'
  return format(parseISO(iso), 'HH:mm')
}

/** Showtime in Astana / Almaty (UTC+5). */
export function formatTimeLocal(iso: string): string {
  if (!iso) return '—'
  return parseISO(iso).toLocaleTimeString('en-US', {
    timeZone: ALMATY_TZ,
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  })
}

export function formatDateLocal(iso: string): string {
  if (!iso) return '—'
  return parseISO(iso).toLocaleDateString('en-GB', {
    timeZone: ALMATY_TZ,
    weekday: 'short',
    day: 'numeric',
    month: 'short',
  })
}

/** YYYY-MM-DD in Asia/Almaty for comparing session days. */
export function localDayKey(d: Date): string {
  return new Intl.DateTimeFormat('en-CA', { timeZone: ALMATY_TZ }).format(d)
}

export function sessionLocalDayKey(iso: string): string {
  return new Intl.DateTimeFormat('en-CA', { timeZone: ALMATY_TZ }).format(parseISO(iso))
}

export function countdown(iso: string): string {
  if (!iso) return '—'
  const d = parseISO(iso)
  if (!isFuture(d)) return 'Started'
  return formatDistanceToNow(d, { addSuffix: true })
}
