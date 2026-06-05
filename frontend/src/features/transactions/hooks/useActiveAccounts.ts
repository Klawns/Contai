import { useQuery } from '@tanstack/react-query'
import { listActiveAccounts } from '../services/accountService.ts'
import { accountsQueryKeys } from './queryKeys.ts'

export function useActiveAccounts() {
  return useQuery({
    queryKey: accountsQueryKeys.all,
    queryFn: listActiveAccounts,
  })
}
