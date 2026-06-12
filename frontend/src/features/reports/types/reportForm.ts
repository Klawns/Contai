import type {
  ReportGroupBy,
  ReportMovementType,
  ReportSettlementStatus,
} from './reports.ts'

export type ReportOption<TValue extends string> = {
  value: TValue
  label: string
}

export type ReportFormState = {
  startDate: string
  endDate: string
  movementType: ReportMovementType
  categoryId: string
  accountId: string
  settlementStatus: ReportSettlementStatus
  groupBy: ReportGroupBy
}
