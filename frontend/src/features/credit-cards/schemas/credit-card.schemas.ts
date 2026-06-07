import { z } from 'zod'
import { creditCardStatuses } from '../types/credit-card.types.ts'
import { rfc3339DateTimeSchema } from './shared.schemas.ts'

export const creditCardSchema = z.object({
  id: z.string(),
  userId: z.string(),
  name: z.string(),
  linkedAccountId: z.string(),
  limitTotal: z.number().int(),
  limitUsed: z.number().int(),
  limitAvailable: z.number().int(),
  closingDay: z.number().int(),
  dueDay: z.number().int(),
  status: z.enum(creditCardStatuses),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const creditCardsSchema = z.array(creditCardSchema)

export const creditCardPayloadSchema = z.object({
  name: z.string().trim().min(1, 'Informe o nome.'),
  linkedAccountId: z.string().trim().min(1, 'Selecione a conta vinculada.'),
  limitTotal: z.number().int().positive('Informe um limite maior que zero.'),
  closingDay: z.number().int().min(1, 'Informe um dia valido.').max(31, 'Informe um dia valido.'),
  dueDay: z.number().int().min(1, 'Informe um dia valido.').max(31, 'Informe um dia valido.'),
  status: z.enum(creditCardStatuses).optional(),
})

export const cardFormSchema = creditCardPayloadSchema.extend({
  status: z.enum(creditCardStatuses),
})
