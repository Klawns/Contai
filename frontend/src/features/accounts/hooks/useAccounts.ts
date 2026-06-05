import { useQuery } from '@tanstack/react-query'
import { getTotalBalance, listActiveAccounts } from '../services/accountService.ts'
import { accountQueryKeys } from './queryKeys.ts'

export function useAccounts() {
  return useQuery({
    queryKey: accountQueryKeys.list,
    queryFn: listActiveAccounts,
  })
}

export function useTotalBalance() {
  return useQuery({
    queryKey: accountQueryKeys.totalBalance,
    queryFn: getTotalBalance,
  })
}
