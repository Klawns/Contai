import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createCategory } from '../services/categoryService.ts'
import type {
  CategoryTransactionType,
  CreateCategoryPayload,
} from '../types/transactions.ts'
import { categoriesQueryKeys } from './queryKeys.ts'

export function useCreateCategory(type: CategoryTransactionType) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: Omit<CreateCategoryPayload, 'type'>) =>
      createCategory({ ...payload, type }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: categoriesQueryKeys.list(type) })
    },
  })
}
