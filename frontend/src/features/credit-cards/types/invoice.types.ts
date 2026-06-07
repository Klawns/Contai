import type { CardPurchaseStatus } from './purchase.types.ts'

export const cardInvoiceStatuses = ['open', 'closed', 'paid', 'canceled'] as const
export const cardInvoiceEffectiveStatuses = ['open', 'closed', 'overdue', 'paid', 'canceled'] as const

export type CardInvoiceStatus = (typeof cardInvoiceStatuses)[number]
export type CardInvoiceEffectiveStatus = (typeof cardInvoiceEffectiveStatuses)[number]

export type CardInstallment = {
  id: string
  userId: string
  cardId: string
  purchaseId: string
  invoiceId: string
  number: number
  amount: number
  status: CardPurchaseStatus
  referenceMonth: string
  createdAt: string
  updatedAt: string
}

export type CardInvoice = {
  id: string
  userId: string
  cardId: string
  referenceMonth: string
  closingAt: string
  dueAt: string
  amount: number
  status: CardInvoiceStatus
  effectiveStatus: CardInvoiceEffectiveStatus
  paidAt: string | null
  paymentTransactionId: string | null
  installments: CardInstallment[]
  createdAt: string
  updatedAt: string
}

export type PayInvoicePayload = {
  occurredAt: string
  categoryId: string
  note: string
}
