export const reportMovementTypes = [
  'all',
  'income',
  'expense',
  'credit_card_expense',
  'transfer',
] as const

export const reportSettlementStatuses = ['all', 'settled', 'pending'] as const
export const reportGroupByOptions = ['none', 'category', 'account', 'day', 'month'] as const

export type ReportMovementType = (typeof reportMovementTypes)[number]
export type ReportSettlementStatus = (typeof reportSettlementStatuses)[number]
export type ReportGroupBy = (typeof reportGroupByOptions)[number]

export type FinancialReportFilters = {
  startAt: string
  endAt: string
  movementType: ReportMovementType
  categoryId?: string
  accountId?: string
  settlementStatus: ReportSettlementStatus
  groupBy: ReportGroupBy
}

export type FinancialReportSummary = {
  incomeTotal: number
  expenseTotal: number
  periodResult: number
  pendingTotal: number
  settledTotal: number
}

export type FinancialMovement = {
  id: string
  source: ReportMovementType
  type: ReportMovementType
  description: string
  amount: number
  occurredAt: string
  categoryId: string | null
  categoryName: string
  accountId: string | null
  accountName: string
  settlementStatus: Exclude<ReportSettlementStatus, 'all'>
}

export type FinancialReportGroup = {
  key: string
  label: string
  incomeTotal: number
  expenseTotal: number
  netTotal: number
  total: number
  count: number
}

export type FinancialReportSeriesPoint = {
  key: string
  label: string
  incomeTotal: number
  expenseTotal: number
  netTotal: number
}

export type FinancialReportCategoryChart = {
  categoryId: string
  name: string
  total: number
}

export type FinancialReportCharts = {
  incomeVsExpense: FinancialReportSeriesPoint[]
  expensesByCategory: FinancialReportCategoryChart[]
  evolution: FinancialReportSeriesPoint[]
}

export type FinancialReport = {
  summary: FinancialReportSummary
  movements: FinancialMovement[]
  groups: FinancialReportGroup[]
  charts: FinancialReportCharts
}
