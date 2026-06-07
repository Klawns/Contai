import { z } from 'zod'
import { cardInvoiceEffectiveStatuses, cardInvoiceStatuses } from '../types/invoice.types.ts'
import { cardPurchaseStatuses } from '../types/purchase.types.ts'
import { rfc3339DateTimeSchema } from './shared.schemas.ts'

export const cardInstallmentSchema = z.object({
  id: z.string(),
  userId: z.string(),
  cardId: z.string(),
  purchaseId: z.string(),
  invoiceId: z.string(),
  number: z.number().int(),
  amount: z.number().int(),
  status: z.enum(cardPurchaseStatuses),
  referenceMonth: rfc3339DateTimeSchema,
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const cardInvoiceSchema = z.object({
  id: z.string(),
  userId: z.string(),
  cardId: z.string(),
  referenceMonth: rfc3339DateTimeSchema,
  closingAt: rfc3339DateTimeSchema,
  dueAt: rfc3339DateTimeSchema,
  amount: z.number().int(),
  status: z.enum(cardInvoiceStatuses),
  effectiveStatus: z.enum(cardInvoiceEffectiveStatuses),
  paidAt: rfc3339DateTimeSchema.nullable(),
  paymentTransactionId: z.string().nullable(),
  installments: z.array(cardInstallmentSchema),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const cardInvoicesSchema = z.array(cardInvoiceSchema)

export const payInvoicePayloadSchema = z.object({
  occurredAt: rfc3339DateTimeSchema,
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  note: z.string(),
})

export const payInvoiceFormSchema = z.object({
  occurredOn: z.string().min(1, 'Informe a data.'),
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  note: z.string(),
})
