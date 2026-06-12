import type { Category } from '../../transactions/types/transactions.ts'
import { formatLocalRFC3339, fromDateInputValue, toDateInputValue } from '../../transactions/utils/date.ts'
import type { ReportFormState } from '../types/reportForm.ts'
import type { FinancialReportFilters } from '../types/reports.ts'

export const allReportOption = 'all'

export function getDefaultReportFormState(): ReportFormState {
  const now = new Date()

  return {
    startDate: toDateInputValue(new Date(now.getFullYear(), now.getMonth(), 1)),
    endDate: toDateInputValue(new Date(now.getFullYear(), now.getMonth() + 1, 0)),
    movementType: 'all',
    categoryId: allReportOption,
    accountId: allReportOption,
    settlementStatus: 'all',
    groupBy: 'none',
  }
}

export function buildFinancialReportFilters(state: ReportFormState): FinancialReportFilters {
  const startDate = fromDateInputValue(state.startDate)
  const endDate = fromDateInputValue(state.endDate)

  startDate.setHours(0, 0, 0, 0)
  endDate.setHours(23, 59, 59, 0)

  return {
    startAt: formatLocalRFC3339(startDate),
    endAt: formatLocalRFC3339(endDate),
    movementType: state.movementType,
    settlementStatus: state.settlementStatus,
    groupBy: state.groupBy,
    ...(state.categoryId === allReportOption ? {} : { categoryId: state.categoryId }),
    ...(state.accountId === allReportOption ? {} : { accountId: state.accountId }),
  }
}

export function formatReportDateLabel(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(fromDateInputValue(value))
}

export function uniqueReportCategories(categories: Category[]) {
  const byId = new Map<string, Category>()

  categories.forEach((category) => byId.set(category.id, category))

  return Array.from(byId.values()).sort((first, second) => first.name.localeCompare(second.name))
}
