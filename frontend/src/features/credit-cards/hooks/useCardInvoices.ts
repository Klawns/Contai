import { useQuery } from '@tanstack/react-query'
import { getCardInvoice, listCardInvoices } from '../services/creditCardService.ts'
import { creditCardQueryKeys } from './queryKeys.ts'

export function useCardInvoices(cardId: string) {
  return useQuery({
    queryKey: creditCardQueryKeys.invoices(cardId),
    queryFn: () => listCardInvoices(cardId),
    enabled: Boolean(cardId),
  })
}

export function useCardInvoice(invoiceId: string) {
  return useQuery({
    queryKey: creditCardQueryKeys.invoice(invoiceId),
    queryFn: () => getCardInvoice(invoiceId),
    enabled: Boolean(invoiceId),
  })
}
