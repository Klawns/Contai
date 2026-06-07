import { useQueryClient } from '@tanstack/react-query'
import { accountQueryKeys } from '../../accounts/hooks/queryKeys.ts'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { accountsQueryKeys, transactionsQueryKeys } from '../../transactions/hooks/queryKeys.ts'
import { creditCardQueryKeys } from './queryKeys.ts'

export function useInvalidateCreditCardData() {
  const queryClient = useQueryClient()

  return () => {
    void queryClient.invalidateQueries({ queryKey: creditCardQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
  }
}

export function useInvalidateFinancialData() {
  const queryClient = useQueryClient()

  return () => {
    void queryClient.invalidateQueries({ queryKey: creditCardQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: transactionsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: accountQueryKeys.all })
  }
}
