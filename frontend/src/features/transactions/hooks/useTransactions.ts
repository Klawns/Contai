import { useQuery } from '@tanstack/react-query'
import { listTransactions } from '../services/transactionService.ts'
import type { TransactionFilters } from '../types/transactions.ts'
import { transactionsQueryKeys } from './queryKeys.ts'

export function useTransactions(filters: TransactionFilters) {
  return useQuery({
    queryKey: transactionsQueryKeys.list(filters),
    queryFn: () => listTransactions(filters),
  })
}
