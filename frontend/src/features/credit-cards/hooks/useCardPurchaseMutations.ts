import { useMutation } from '@tanstack/react-query'
import { cancelCardPurchase, createCardPurchase } from '../services/creditCardService.ts'
import type { CardPurchasePayload } from '../types/purchase.types.ts'
import { useInvalidateCreditCardData } from './useCreditCardMutationInvalidation.ts'

export function useCreateCardPurchase(cardId: string) {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: (payload: CardPurchasePayload) => createCardPurchase(cardId, payload),
    onSuccess: invalidate,
  })
}

export function useCancelCardPurchase() {
  const invalidate = useInvalidateCreditCardData()

  return useMutation({
    mutationFn: cancelCardPurchase,
    onSuccess: invalidate,
  })
}
