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

export function formatTime(iso: string): string {
  if (!iso) return '—'
  return format(parseISO(iso), 'HH:mm')
}

export function countdown(iso: string): string {
  if (!iso) return '—'
  const d = parseISO(iso)
  if (!isFuture(d)) return 'Started'
  return formatDistanceToNow(d, { addSuffix: true })
}
