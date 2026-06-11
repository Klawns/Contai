import type { CreditCardStatus } from './credit-card.types.ts'
import type { CardPurchaseType } from './purchase.types.ts'

export type CardFormValues = {
  name: string
  linkedAccountId: string
  limitTotal: number
  closingDay: number
  dueDay: number
  status: CreditCardStatus
}

export type PurchaseFormValues = {
  cardId: string
  categoryId: string
  description: string
  totalAmount: number
  purchaseDate: string
  purchaseType: CardPurchaseType
  installmentCount: number
  firstInvoiceMonth: string
  note: string
}

export type PayInvoiceFormValues = {
  occurredOn: string
  categoryId: string
  note: string
}
