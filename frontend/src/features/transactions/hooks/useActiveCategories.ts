import { useQuery } from '@tanstack/react-query'
import { listActiveCategories } from '../services/categoryService.ts'
import type { CategoryTransactionType } from '../types/transactions.ts'
import { categoriesQueryKeys } from './queryKeys.ts'

export function useActiveCategories(type: CategoryTransactionType) {
  return useQuery({
    queryKey: categoriesQueryKeys.list(type),
    queryFn: () => listActiveCategories(type),
  })
}
