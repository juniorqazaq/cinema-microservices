export type TicketCategoryId =
  | 'child'
  | 'student'
  | 'teen'
  | 'adult'
  | 'senior'

export interface TicketCategory {
  id: TicketCategoryId
  label: string
  discount: number
  description: string
}

export const TICKET_CATEGORIES: TicketCategory[] = [
  { id: 'child', label: 'Child (under 12)', discount: 0.5, description: '-50%' },
  { id: 'student', label: 'Student', discount: 0.3, description: '-30%' },
  { id: 'teen', label: 'Teen (12–17)', discount: 0.2, description: '-20%' },
  { id: 'adult', label: 'Adult (18+)', discount: 0, description: 'Full price' },
  { id: 'senior', label: 'Senior (60+)', discount: 0.4, description: '-40%' },
]

export function ticketPrice(basePrice: number, categoryId: TicketCategoryId): number {
  const cat = TICKET_CATEGORIES.find((c) => c.id === categoryId)
  if (!cat) return basePrice
  return Math.round(basePrice * (1 - cat.discount))
}

export function categoryLabel(id: string): string {
  return TICKET_CATEGORIES.find((c) => c.id === id)?.label ?? id
}
