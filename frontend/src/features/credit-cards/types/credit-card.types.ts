export const creditCardStatuses = ['active', 'inactive'] as const

export type CreditCardStatus = (typeof creditCardStatuses)[number]

export type CreditCard = {
  id: string
  userId: string
  name: string
  linkedAccountId: string
  limitTotal: number
  limitUsed: number
  limitAvailable: number
  closingDay: number
  dueDay: number
  status: CreditCardStatus
  createdAt: string
  updatedAt: string
}

export type CreditCardPayload = {
  name: string
  linkedAccountId: string
  limitTotal: number
  closingDay: number
  dueDay: number
  status?: CreditCardStatus
}
