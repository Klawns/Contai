import { useQuery } from '@tanstack/react-query'
import { listActiveAccounts } from '../../transactions/services/accountService.ts'
import { reportsQueryKeys } from './queryKeys.ts'

export function useReportAccounts() {
  return useQuery({
    queryKey: reportsQueryKeys.accounts(),
    queryFn: listActiveAccounts,
  })
}
