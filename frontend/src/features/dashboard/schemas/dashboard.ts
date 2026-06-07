import { z } from 'zod'

const rfc3339Pattern =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/

export const rfc3339DateTimeSchema = z
  .string()
  .regex(rfc3339Pattern, 'Use uma data RFC3339 valida.')
  .refine((value) => !Number.isNaN(Date.parse(value)), {
    message: 'Use uma data RFC3339 valida.',
  })

export const dashboardPeriodSchema = z
  .object({
    startAt: rfc3339DateTimeSchema,
    endAt: rfc3339DateTimeSchema,
  })
  .refine(({ startAt, endAt }) => Date.parse(startAt) <= Date.parse(endAt), {
    message: 'A data inicial deve ser anterior ou igual a data final.',
    path: ['endAt'],
  })

export const monthlyDashboardFiltersSchema = dashboardPeriodSchema

export const accountBalanceSchema = z.object({
  accountId: z.string(),
  name: z.string(),
  type: z.string(),
  balance: z.number().int(),
  bankIconId: z.string(),
})

export const creditCardDashboardSchema = z.object({
  cardId: z.string(),
  name: z.string(),
  linkedAccountId: z.string(),
  limitTotal: z.number().int(),
  limitUsed: z.number().int(),
  limitAvailable: z.number().int(),
  currentInvoiceId: z.string().nullable(),
  currentInvoiceAmount: z.number().int(),
  currentInvoiceDueAt: rfc3339DateTimeSchema.nullable(),
  currentInvoiceEffectiveStatus: z.string(),
})

export const expenseByCategorySchema = z.object({
  categoryId: z.string(),
  name: z.string(),
  color: z.string(),
  icon: z.string(),
  total: z.number().int(),
})

export const recentTransactionSchema = z.object({
  id: z.string(),
  userId: z.string(),
  type: z.string(),
  description: z.string(),
  amount: z.number().int(),
  occurredAt: rfc3339DateTimeSchema,
  accountId: z.string().nullable(),
  sourceAccountId: z.string().nullable(),
  destinationAccountId: z.string().nullable(),
  categoryId: z.string().nullable(),
  status: z.string(),
  note: z.string(),
  removedAt: rfc3339DateTimeSchema.nullable(),
  createdAt: rfc3339DateTimeSchema,
  updatedAt: rfc3339DateTimeSchema,
})

export const monthlyDashboardSchema = z.object({
  userId: z.string(),
  period: dashboardPeriodSchema,
  totalBalance: z.number().int(),
  monthlyIncome: z.number().int(),
  monthlyExpense: z.number().int(),
  monthlyTransferIn: z.number().int(),
  monthlyTransferOut: z.number().int(),
  monthlyNetBalance: z.number().int(),
  accountBalances: z.array(accountBalanceSchema),
  creditCards: z.array(creditCardDashboardSchema),
  expensesByCategory: z.array(expenseByCategorySchema),
  recentTransactions: z.array(recentTransactionSchema),
})

export const monthlySeriesPointSchema = z.object({
  monthStartAt: rfc3339DateTimeSchema,
  monthEndAt: rfc3339DateTimeSchema,
  income: z.number().int(),
  expense: z.number().int(),
  balance: z.number().int(),
})

export const monthlySeriesSchema = z.object({
  userId: z.string(),
  period: dashboardPeriodSchema,
  points: z.array(monthlySeriesPointSchema),
})
