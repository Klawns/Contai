import { formatLocalRFC3339, fromDateInputValue } from '../../transactions/utils/date.ts'

export function formatCardDate(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(new Date(value))
}

export function formatInvoiceMonth(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    month: 'long',
    year: 'numeric',
  }).format(new Date(value))
}

export function toLocalDateTime(value: string) {
  return formatLocalRFC3339(fromDateInputValue(value))
}
