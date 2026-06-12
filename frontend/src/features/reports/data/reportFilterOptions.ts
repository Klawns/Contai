import type { ReportOption } from '../types/reportForm.ts'
import type {
  ReportGroupBy,
  ReportMovementType,
  ReportSettlementStatus,
} from '../types/reports.ts'

export const movementTypeOptions = [
  { value: 'all', label: 'Todos' },
  { value: 'income', label: 'Receita' },
  { value: 'expense', label: 'Despesa' },
  { value: 'credit_card_expense', label: 'Despesa cartao' },
  { value: 'transfer', label: 'Transferencia' },
] satisfies Array<ReportOption<ReportMovementType>>

export const settlementOptions = [
  { value: 'all', label: 'Todos' },
  { value: 'settled', label: 'Pago/Recebido' },
  { value: 'pending', label: 'Nao pago/Nao recebido' },
] satisfies Array<ReportOption<ReportSettlementStatus>>

export const groupOptions = [
  { value: 'none', label: 'Nenhum' },
  { value: 'category', label: 'Por categoria' },
  { value: 'account', label: 'Por conta' },
  { value: 'day', label: 'Por dia' },
  { value: 'month', label: 'Por mes' },
] satisfies Array<ReportOption<ReportGroupBy>>
