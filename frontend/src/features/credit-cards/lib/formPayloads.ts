import type { CreditCardPayload } from '../types/credit-card.types.ts'
import type { CardFormValues, PayInvoiceFormValues, PurchaseFormValues } from '../types/form.types.ts'
import type { PayInvoicePayload } from '../types/invoice.types.ts'
import type { CardPurchasePayload } from '../types/purchase.types.ts'
import { toLocalDateTime } from './creditCardDates.ts'

export function toCreditCardPayload(values: CardFormValues, mode: 'create' | 'edit'): CreditCardPayload {
  return {
    name: values.name.trim(),
    linkedAccountId: values.linkedAccountId,
    limitTotal: values.limitTotal,
    closingDay: values.closingDay,
    dueDay: values.dueDay,
    status: mode === 'edit' ? values.status : undefined,
  }
}

export function toCardPurchasePayload(values: PurchaseFormValues): CardPurchasePayload {
  return {
    categoryId: values.categoryId,
    description: values.description.trim(),
    totalAmount: values.totalAmount,
    purchaseDate: toLocalDateTime(values.purchaseDate),
    purchaseType: values.purchaseType,
    installmentCount: values.installmentCount,
    firstInvoiceMonth: values.firstInvoiceMonth,
    note: values.note.trim(),
  }
}

export function toPayInvoicePayload(values: PayInvoiceFormValues): PayInvoicePayload {
  return {
    occurredAt: toLocalDateTime(values.occurredOn),
    categoryId: values.categoryId,
    note: values.note.trim(),
  }
}
