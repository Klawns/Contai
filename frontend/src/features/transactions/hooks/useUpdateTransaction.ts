import { useMutation, useQueryClient } from '@tanstack/react-query'
import { accountQueryKeys } from '../../accounts/hooks/queryKeys.ts'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { updateTransaction } from '../services/transactionService.ts'
import type {
  TransactionType,
  UpdateTransactionPayloadByType,
} from '../types/transactions.ts'
import { accountsQueryKeys, transactionsQueryKeys } from './queryKeys.ts'

export function useUpdateTransaction<TType extends TransactionType>(
  type: TType,
  transactionId: string,
) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: UpdateTransactionPayloadByType[TType]) =>
      updateTransaction(type, transactionId, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: transactionsQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: accountQueryKeys.all })
    },
  })
}
