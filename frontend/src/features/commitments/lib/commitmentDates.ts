import type { SelectedMonth } from '../../../components/MonthSelector.tsx'
import {
  formatLocalRFC3339,
  fromDateInputValue,
} from '../../transactions/utils/date.ts'

const monthQueryPattern = /^(\d{4})-(0[1-9]|1[0-2])$/

export function parseMonthQuery(value: string | null): SelectedMonth | null {
  const match = value?.match(monthQueryPattern)

  if (!match) {
    return null
  }

  return {
    year: Number(match[1]),
    monthIndex: Number(match[2]) - 1,
  }
}

export function formatCommitmentDate(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(new Date(value))
}

export function toLocalDateTime(value: string) {
  return formatLocalRFC3339(fromDateInputValue(value))
}
