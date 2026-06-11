import { z } from 'zod'
import { cardPurchaseStatuses, cardPurchaseTypes } from '../types/purchase.types.ts'
import { rfc3339DateTimeSchema } from './shared.schemas.ts'

export const cardPurchaseSchema = z.object({
  id: z.string(),
  userId: z.string(),
  cardId: z.string(),
  categoryId: z.string(),
  description: z.string(),
  totalAmount: z.number().int(),
  purchaseDate: rfc3339DateTimeSchema,
  purchaseType: z.enum(cardPurchaseTypes).default('single'),
  installmentCount: z.number().int(),
  firstInvoiceMonth: z.string().regex(/^\d{4}-\d{2}$/).default('1970-01'),
  note: z.string(),
  status: z.enum(cardPurchaseStatuses),
  canceledAt: rfc3339DateTimeSchema.nullable(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const cardPurchasesSchema = z.array(cardPurchaseSchema)

export const cardPurchasePayloadSchema = z.object({
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  description: z.string().trim().min(1, 'Informe a descricao.'),
  totalAmount: z.number().int().positive('Informe um valor maior que zero.'),
  purchaseDate: rfc3339DateTimeSchema,
  purchaseType: z.enum(cardPurchaseTypes),
  installmentCount: z.number().int().min(1, 'Informe ao menos uma parcela.').max(12, 'Use ate 12 parcelas.'),
  firstInvoiceMonth: z.string().regex(/^\d{4}-\d{2}$/, 'Use o formato AAAA-MM.'),
  note: z.string(),
})

export const purchaseFormSchema = z.object({
  cardId: z.string().trim().min(1, 'Selecione um cartao.'),
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  description: z.string().trim().min(1, 'Informe a descricao.'),
  totalAmount: z.number().int().positive('Informe um valor maior que zero.'),
  purchaseDate: z.string().min(1, 'Informe a data.'),
  purchaseType: z.enum(cardPurchaseTypes),
  installmentCount: z.number().int().min(1, 'Informe ao menos uma parcela.').max(12, 'Use ate 12 parcelas.'),
  firstInvoiceMonth: z.string().regex(/^\d{4}-\d{2}$/, 'Use o formato AAAA-MM.'),
  note: z.string(),
}).superRefine((values, context) => {
  if (values.purchaseType === 'single' && values.installmentCount !== 1) {
    context.addIssue({ code: 'custom', path: ['installmentCount'], message: 'Compra unica usa 1x.' })
  }
  if (values.purchaseType === 'installment' && values.installmentCount < 2) {
    context.addIssue({ code: 'custom', path: ['installmentCount'], message: 'Parcelada comeca em 2x.' })
  }
  if (values.purchaseType === 'fixed' && values.installmentCount !== 1) {
    context.addIssue({ code: 'custom', path: ['installmentCount'], message: 'Compra fixa nao usa parcelamento.' })
  }
})
