import { z } from 'zod'
import {
  categoryTransactionTypes,
  transactionOriginTypes,
  transactionTypes,
  type CategoryTransactionType,
} from '../types/transactions.ts'

const rfc3339Pattern =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/

export const rfc3339DateTimeSchema = z
  .string()
  .regex(rfc3339Pattern, 'Use uma data RFC3339 valida.')
  .refine((value) => !Number.isNaN(Date.parse(value)), {
    message: 'Use uma data RFC3339 valida.',
  })

export const transactionTypeSchema = z.enum(transactionTypes)
export const transactionOriginTypeSchema = z.enum(transactionOriginTypes)
export const categoryTransactionTypeSchema = z.enum(categoryTransactionTypes)

export const transactionFiltersSchema = z
  .object({
    startAt: rfc3339DateTimeSchema.optional(),
    endAt: rfc3339DateTimeSchema.optional(),
    accountId: z.string().min(1).optional(),
    categoryId: z.string().min(1).optional(),
    type: transactionTypeSchema.optional(),
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

export const transactionSchema = z.object({
  id: z.string(),
  userId: z.string(),
  type: transactionTypeSchema,
  description: z.string(),
  amount: z.number().int(),
  occurredAt: rfc3339DateTimeSchema,
  accountId: z.string().nullable(),
  sourceAccountId: z.string().nullable(),
  destinationAccountId: z.string().nullable(),
  categoryId: z.string().nullable(),
  status: z.string(),
  originType: transactionOriginTypeSchema.default('manual'),
  originId: z.string().nullable().default(null),
  note: z.string(),
  removedAt: rfc3339DateTimeSchema.nullable(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const transactionsSchema = z.array(transactionSchema)

export const accountSchema = z.object({
  id: z.string(),
  userId: z.string(),
  name: z.string(),
  type: z.string(),
  initialBalance: z.number().int(),
  currentBalance: z.number().int(),
  bankIconId: z.string(),
  includeInDashboardTotal: z.boolean(),
  status: z.string(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const accountsSchema = z.array(accountSchema)

export const categorySchema = z.object({
  id: z.string(),
  userId: z.string(),
  name: z.string(),
  normalizedName: z.string(),
  type: categoryTransactionTypeSchema,
  color: z.string(),
  icon: z.string(),
  isDefault: z.boolean(),
  status: z.string(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const categoriesSchema = z.array(categorySchema)

const createBaseTransactionPayloadSchema = z.object({
  description: z.string().trim().min(1, 'Informe a descricao.'),
  amount: z.number().int().positive('Informe um valor maior que zero.'),
  occurredAt: rfc3339DateTimeSchema,
  note: z.string(),
})

export const createIncomeExpenseTransactionPayloadSchema =
  createBaseTransactionPayloadSchema.extend({
    accountId: z.string().trim().min(1, 'Selecione uma conta.'),
    categoryId: z.string().trim().min(1, 'Selecione uma categoria.'),
  })

export const createTransferTransactionPayloadSchema =
  createBaseTransactionPayloadSchema
    .extend({
      sourceAccountId: z.string().trim().min(1, 'Selecione a conta de origem.'),
      destinationAccountId: z
        .string()
        .trim()
        .min(1, 'Selecione a conta de destino.'),
    })
    .refine(
      ({ sourceAccountId, destinationAccountId }) =>
        sourceAccountId !== destinationAccountId,
      {
        message: 'Origem e destino devem ser diferentes.',
        path: ['destinationAccountId'],
      },
    )

export const updateIncomeExpenseTransactionPayloadSchema =
  createIncomeExpenseTransactionPayloadSchema

export const updateTransferTransactionPayloadSchema =
  createTransferTransactionPayloadSchema

export const createCategoryPayloadSchema = z.object({
  name: z.string().trim().min(1, 'Informe o nome.'),
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Use uma cor hexadecimal.'),
  icon: z.string().trim().min(1, 'Informe o icone.'),
  type: z.enum(categoryTransactionTypes) satisfies z.ZodType<CategoryTransactionType>,
})
