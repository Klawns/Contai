import { useMutation, useQueryClient } from '@tanstack/react-query'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import {
  createExpenseTransaction,
  createIncomeTransaction,
  createTransferTransaction,
} from '../services/transactionService.ts'
import type {
  CreateTransactionPayloadByType,
  TransactionType,
} from '../types/transactions.ts'
import { accountsQueryKeys, transactionsQueryKeys } from './queryKeys.ts'

const mutationByType = {
  income: createIncomeTransaction,
  expense: createExpenseTransaction,
  transfer: createTransferTransaction,
}

export function useCreateTransaction<TType extends TransactionType>(type: TType) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateTransactionPayloadByType[TType]) =>
      mutationByType[type](payload as never),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: transactionsQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
      void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
    },
  })
}
