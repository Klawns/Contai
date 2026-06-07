import { z } from 'zod'
import { cardPurchaseStatuses } from '../types/purchase.types.ts'
import { rfc3339DateTimeSchema } from './shared.schemas.ts'

export const cardPurchaseSchema = z.object({
  id: z.string(),
  userId: z.string(),
  cardId: z.string(),
  categoryId: z.string(),
  description: z.string(),
  totalAmount: z.number().int(),
  purchaseDate: rfc3339DateTimeSchema,
  installmentCount: z.number().int(),
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
  installmentCount: z.number().int().min(1, 'Informe ao menos uma parcela.'),
  note: z.string(),
})

export const purchaseFormSchema = z.object({
  cardId: z.string().trim().min(1, 'Selecione um cartao.'),
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  description: z.string().trim().min(1, 'Informe a descricao.'),
  totalAmount: z.number().int().positive('Informe um valor maior que zero.'),
  purchaseDate: z.string().min(1, 'Informe a data.'),
  installmentCount: z.number().int().min(1, 'Informe ao menos uma parcela.').max(48, 'Use ate 48 parcelas.'),
  note: z.string(),
})
