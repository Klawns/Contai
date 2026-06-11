import { formatCurrency } from '../../transactions/utils/money.ts'

export function getInstallmentText(amount: number, installmentCount: number) {
  return `${installmentCount}x / ${formatCurrency(Math.round(amount / installmentCount))}`
}
