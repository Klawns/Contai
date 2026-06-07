import { useQuery } from '@tanstack/react-query'
import { listCreditCards } from '../services/creditCardService.ts'
import { creditCardQueryKeys } from './queryKeys.ts'

export function useCreditCards() {
  return useQuery({
    queryKey: creditCardQueryKeys.lists(),
    queryFn: listCreditCards,
  })
}
