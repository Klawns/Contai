import { useMutation, useQueryClient } from '@tanstack/react-query'
import { accountQueryKeys } from '../../accounts/hooks/queryKeys.ts'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { deleteTransaction } from '../services/transactionService.ts'
import { accountsQueryKeys, transactionsQueryKeys } from './queryKeys.ts'

export function useDeleteTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (transactionId: string) => deleteTransaction(transactionId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: transactionsQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: accountQueryKeys.all })
    },
  })
}
