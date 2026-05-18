export const CINEMA_ADDRESSES: Record<string, string> = {
  'Keruen Cineplex': 'Keruen Mall, Astana',
  'Khan Shatyr IMAX': 'Khan Shatyr, Astana',
  'Kinopark 8': 'Dostyk Ave, Almaty',
}

export function cinemaAddress(name: string, city: string): string {
  return CINEMA_ADDRESSES[name] ?? `${name}, ${city}`
}
