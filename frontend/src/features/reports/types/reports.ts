export type ReportTransactionType = 'income' | 'expense'

export type ReportPeriodFilters = {
  startAt: string
  endAt: string
}

export type TransactionsReportFilters = ReportPeriodFilters & {
  type: ReportTransactionType
}

export type AccountReportFilters = ReportPeriodFilters & {
  accountId: string
}
