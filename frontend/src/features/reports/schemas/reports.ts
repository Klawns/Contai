import { z } from 'zod'
import {
  reportGroupByOptions,
  reportMovementTypes,
  reportSettlementStatuses,
} from '../types/reports.ts'

const rfc3339Pattern =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/

export const rfc3339DateTimeSchema = z
  .string()
  .regex(rfc3339Pattern, 'Use uma data RFC3339 valida.')
  .refine((value) => !Number.isNaN(Date.parse(value)), {
    message: 'Use uma data RFC3339 valida.',
  })

export const financialReportFiltersSchema = z
  .object({
    startAt: rfc3339DateTimeSchema,
    endAt: rfc3339DateTimeSchema,
    movementType: z.enum(reportMovementTypes),
    categoryId: z.string().min(1).optional(),
    accountId: z.string().min(1).optional(),
    settlementStatus: z.enum(reportSettlementStatuses),
    groupBy: z.enum(reportGroupByOptions),
  })
  .refine(({ startAt, endAt }) => Date.parse(startAt) <= Date.parse(endAt), {
    message: 'A data inicial deve ser anterior ou igual a data final.',
    path: ['endAt'],
  })

const financialReportSummarySchema = z.object({
  incomeTotal: z.number().int(),
  expenseTotal: z.number().int(),
  periodResult: z.number().int(),
  pendingTotal: z.number().int(),
  settledTotal: z.number().int(),
})

const financialMovementSchema = z.object({
  id: z.string(),
  source: z.enum(reportMovementTypes),
  type: z.enum(reportMovementTypes),
  description: z.string(),
  amount: z.number().int(),
  occurredAt: rfc3339DateTimeSchema,
  categoryId: z.string().nullable(),
  categoryName: z.string(),
  accountId: z.string().nullable(),
  accountName: z.string(),
  settlementStatus: z.enum(['settled', 'pending']),
})

const financialReportGroupSchema = z.object({
  key: z.string(),
  label: z.string(),
  incomeTotal: z.number().int(),
  expenseTotal: z.number().int(),
  netTotal: z.number().int(),
  total: z.number().int(),
  count: z.number().int(),
})

const financialReportSeriesPointSchema = z.object({
  key: z.string(),
  label: z.string(),
  incomeTotal: z.number().int(),
  expenseTotal: z.number().int(),
  netTotal: z.number().int(),
})

const financialReportCategoryChartSchema = z.object({
  categoryId: z.string(),
  name: z.string(),
  total: z.number().int(),
})

export const financialReportSchema = z.object({
  summary: financialReportSummarySchema,
  movements: z.array(financialMovementSchema),
  groups: z.array(financialReportGroupSchema),
  charts: z.object({
    incomeVsExpense: z.array(financialReportSeriesPointSchema),
    expensesByCategory: z.array(financialReportCategoryChartSchema),
    evolution: z.array(financialReportSeriesPointSchema),
  }),
})
