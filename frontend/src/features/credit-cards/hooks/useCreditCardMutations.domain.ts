import { useMutation } from '@tanstack/react-query'
import {
  createCreditCard,
  inactivateCreditCard,
  updateCreditCard,
} from '../services/creditCardService.ts'
import type { CreditCardPayload } from '../types/credit-card.types.ts'
import { useInvalidateCreditCardData } from './useCreditCardMutationInvalidation.ts'

export function useCreateCreditCard() {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: createCreditCard,
    onSuccess: invalidate,
  })
}

export function useUpdateCreditCard(cardId: string) {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: (payload: CreditCardPayload) => updateCreditCard(cardId, payload),
    onSuccess: invalidate,
  })
}

export function useInactivateCreditCard() {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: inactivateCreditCard,
    onSuccess: invalidate,
  })
}
