import type { TransactionSelectOption } from '../../transactions/components/Selectors.tsx'
import type { CreditCardStatus } from '../types/credit-card.types.ts'

export const cardStatusCopy = {
  active: 'Ativo',
  inactive: 'Inativo',
} as const

export const cardStatusOptions: Array<TransactionSelectOption<CreditCardStatus>> = [
  { value: 'active', label: cardStatusCopy.active },
  { value: 'inactive', label: cardStatusCopy.inactive },
]
