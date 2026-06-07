import { useMutation } from '@tanstack/react-query'
import { closeCardInvoice, payCardInvoice } from '../services/creditCardService.ts'
import type { PayInvoicePayload } from '../types/invoice.types.ts'
import {
  useInvalidateCreditCardData,
  useInvalidateFinancialData,
} from './useCreditCardMutationInvalidation.ts'

export function useCloseCardInvoice() {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: closeCardInvoice,
    onSuccess: invalidate,
  })
}

export function usePayCardInvoice(invoiceId: string) {
  const invalidate = useInvalidateFinancialData()

  return useMutation({
    mutationFn: (payload: PayInvoicePayload) => payCardInvoice(invoiceId, payload),
    onSuccess: invalidate,
  })
}
