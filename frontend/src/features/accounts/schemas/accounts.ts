import { z } from 'zod'
import { rfc3339DateTimeSchema } from '../../transactions/schemas/transactions.ts'
import { accountTypes } from '../types/accounts.ts'

export const bankIconIdSchema = z
  .string()
  .trim()
  .regex(/^[A-Za-z0-9_-]{1,64}$/, 'Selecione um banco.')

export const accountTypeSchema = z.enum(accountTypes)

export const accountSchema = z.object({
  id: z.string(),
  userId: z.string(),
  name: z.string(),
  type: accountTypeSchema,
  initialBalance: z.number().int(),
  currentBalance: z.number().int(),
  bankIconId: z.string(),
  includeInDashboardTotal: z.boolean(),
  status: z.string(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const accountsSchema = z.array(accountSchema)

export const totalBalanceSchema = z.object({
  totalBalance: z.number().int(),
})

export const createAccountPayloadSchema = z.object({
  name: z.string().trim().min(1, 'Informe o nome.'),
  type: accountTypeSchema,
  initialBalance: z.number().int(),
  bankIconId: bankIconIdSchema,
  includeInDashboardTotal: z.boolean(),
})

export const updateAccountPayloadSchema = createAccountPayloadSchema.omit({
  initialBalance: true,
})
