import { api } from '../../../lib/api/axios.ts'
import {
  creditCardPayloadSchema,
  creditCardSchema,
  creditCardsSchema,
} from '../schemas/credit-card.schemas.ts'
import {
  cardInvoiceSchema,
  cardInvoicesSchema,
  payInvoicePayloadSchema,
} from '../schemas/invoice.schemas.ts'
import {
  cardPurchasePayloadSchema,
  cardPurchaseSchema,
  cardPurchasesSchema,
} from '../schemas/purchase.schemas.ts'
import type { CreditCard, CreditCardPayload } from '../types/credit-card.types.ts'
import type { CardInvoice, PayInvoicePayload } from '../types/invoice.types.ts'
import type { CardPurchase, CardPurchasePayload } from '../types/purchase.types.ts'

export async function listCreditCards(): Promise<CreditCard[]> {
  const response = await api.get<unknown>('/credit-cards')

  return creditCardsSchema.parse(response.data)
}

export async function createCreditCard(payload: CreditCardPayload): Promise<CreditCard> {
  const body = creditCardPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/credit-cards', body)

  return creditCardSchema.parse(response.data)
}

export async function updateCreditCard(
  cardId: string,
  payload: CreditCardPayload,
): Promise<CreditCard> {
  const body = creditCardPayloadSchema.parse(payload)
  const response = await api.patch<unknown>(`/credit-cards/${cardId}`, body)

  return creditCardSchema.parse(response.data)
}

export async function inactivateCreditCard(cardId: string): Promise<CreditCard> {
  const response = await api.patch<unknown>(`/credit-cards/${cardId}/inactivate`)

  return creditCardSchema.parse(response.data)
}

export async function listCardPurchases(cardId: string): Promise<CardPurchase[]> {
  const response = await api.get<unknown>(`/credit-cards/${cardId}/purchases`)

  return cardPurchasesSchema.parse(response.data)
}

export async function createCardPurchase(
  cardId: string,
  payload: CardPurchasePayload,
): Promise<CardPurchase> {
  const body = cardPurchasePayloadSchema.parse(payload)
  const response = await api.post<unknown>(`/credit-cards/${cardId}/purchases`, body)

  return cardPurchaseSchema.parse(response.data)
}

export async function cancelCardPurchase(purchaseId: string): Promise<CardPurchase> {
  const response = await api.patch<unknown>(`/credit-card-purchases/${purchaseId}/cancel`)

  return cardPurchaseSchema.parse(response.data)
}

export async function listCardInvoices(cardId: string): Promise<CardInvoice[]> {
  const response = await api.get<unknown>(`/credit-cards/${cardId}/invoices`)

  return cardInvoicesSchema.parse(response.data)
}

export async function getCardInvoice(invoiceId: string): Promise<CardInvoice> {
  const response = await api.get<unknown>(`/credit-card-invoices/${invoiceId}`)

  return cardInvoiceSchema.parse(response.data)
}

export async function closeCardInvoice(invoiceId: string): Promise<CardInvoice> {
  const response = await api.patch<unknown>(`/credit-card-invoices/${invoiceId}/close`)

  return cardInvoiceSchema.parse(response.data)
}

export async function payCardInvoice(
  invoiceId: string,
  payload: PayInvoicePayload,
): Promise<CardInvoice> {
  const body = payInvoicePayloadSchema.parse(payload)
  const response = await api.patch<unknown>(`/credit-card-invoices/${invoiceId}/pay`, body)

  return cardInvoiceSchema.parse(response.data)
}
