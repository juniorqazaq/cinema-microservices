export const CITIES = ['Astana', 'Almaty'] as const
export type CityName = (typeof CITIES)[number]
export const DEFAULT_CITY: CityName = 'Astana'
