import { formatCurrency } from '../../transactions/utils/money.ts'

export function getInstallmentText(amount: number, installmentCount: number) {
  return installmentCount > 1
    ? `${installmentCount}x de ${formatCurrency(Math.round(amount / installmentCount))}`
    : 'A vista'
}
