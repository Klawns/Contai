import { useQuery } from '@tanstack/react-query'
import { listCardPurchases } from '../services/creditCardService.ts'
import { creditCardQueryKeys } from './queryKeys.ts'

export function useCardPurchases(cardId: string) {
  return useQuery({
    queryKey: creditCardQueryKeys.purchases(cardId),
    queryFn: () => listCardPurchases(cardId),
    enabled: Boolean(cardId),
  })
}
