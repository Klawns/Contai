import type { CardInvoiceEffectiveStatus } from '../types/invoice.types.ts'

export const invoiceStatusCopy: Record<CardInvoiceEffectiveStatus, string> = {
  open: 'Aberta',
  closed: 'Fechada',
  overdue: 'Vencida',
  paid: 'Paga',
  canceled: 'Cancelada',
}
