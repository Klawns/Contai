export const cardPurchaseStatuses = ['active', 'canceled'] as const
export const cardPurchaseTypes = ['single', 'installment', 'fixed'] as const

export type CardPurchaseStatus = (typeof cardPurchaseStatuses)[number]
export type CardPurchaseType = (typeof cardPurchaseTypes)[number]

export type CardPurchase = {
  id: string
  userId: string
  cardId: string
  categoryId: string
  description: string
  totalAmount: number
  purchaseDate: string
  purchaseType: CardPurchaseType
  installmentCount: number
  firstInvoiceMonth: string
  note: string
  status: CardPurchaseStatus
  canceledAt: string | null
  createdAt: string
  updatedAt: string
}

export type CardPurchasePayload = {
  categoryId: string
  description: string
  totalAmount: number
  purchaseDate: string
  purchaseType: CardPurchaseType
  installmentCount: number
  firstInvoiceMonth: string
  note: string
}
