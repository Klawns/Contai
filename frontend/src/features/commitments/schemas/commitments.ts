import { z } from 'zod'
import { rfc3339DateTimeSchema } from '../../transactions/schemas/transactions.ts'
import {
  commitmentStatuses,
  commitmentTypes,
  effectiveCommitmentStatuses,
  recurrenceFrequencies,
} from '../types/commitments.ts'

export const commitmentTypeSchema = z.enum(commitmentTypes)
export const commitmentStatusSchema = z.enum(commitmentStatuses)
export const effectiveCommitmentStatusSchema = z.enum(effectiveCommitmentStatuses)
export const recurrenceFrequencySchema = z.enum(recurrenceFrequencies)

export const recurrenceSchema = z.object({
  frequency: recurrenceFrequencySchema,
  interval: z.number().int().positive(),
  endsAt: rfc3339DateTimeSchema.nullable(),
})

export const commitmentSchema = z.object({
  id: z.string(),
  userId: z.string(),
  type: commitmentTypeSchema,
  description: z.string(),
  amount: z.number().int(),
  dueAt: rfc3339DateTimeSchema,
  accountId: z.string(),
  categoryId: z.string(),
  note: z.string(),
  status: commitmentStatusSchema,
  effectiveStatus: effectiveCommitmentStatusSchema,
  recurrence: recurrenceSchema.nullable(),
  settledAt: rfc3339DateTimeSchema.nullable(),
  settlementTransactionId: z.string().nullable(),
  canceledAt: rfc3339DateTimeSchema.nullable(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const commitmentsSchema = z.array(commitmentSchema)

export const commitmentFiltersSchema = z
  .object({
    startAt: rfc3339DateTimeSchema.optional(),
    endAt: rfc3339DateTimeSchema.optional(),
    status: commitmentStatusSchema.optional(),
    effectiveStatus: effectiveCommitmentStatusSchema.optional(),
    accountId: z.string().min(1).optional(),
    categoryId: z.string().min(1).optional(),
    limit: z.number().int().nonnegative().optional(),
    offset: z.number().int().nonnegative().optional(),
  })
  .refine(
    ({ startAt, endAt }) =>
      !startAt || !endAt || Date.parse(startAt) <= Date.parse(endAt),
    {
      message: 'A data inicial deve ser anterior ou igual a data final.',
      path: ['endAt'],
    },
  )

const baseCommitmentPayloadSchema = z.object({
  description: z.string().trim().min(1, 'Informe a descricao.'),
  amount: z.number().int().positive('Informe um valor maior que zero.'),
  dueAt: rfc3339DateTimeSchema,
  accountId: z.string().trim().min(1, 'Selecione uma conta.'),
  categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  note: z.string(),
})

export const commitmentPayloadSchema = baseCommitmentPayloadSchema.extend({
  recurrence: recurrenceSchema.nullable(),
})

export const settlementPayloadSchema = baseCommitmentPayloadSchema
  .omit({ dueAt: true, description: true })
  .extend({
    occurredAt: rfc3339DateTimeSchema,
  })
