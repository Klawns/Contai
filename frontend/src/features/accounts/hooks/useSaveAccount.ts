import { useMutation, useQueryClient } from '@tanstack/react-query'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { accountsQueryKeys } from '../../transactions/hooks/queryKeys.ts'
import { createAccount, deleteAccount, updateAccount } from '../services/accountService.ts'
import type { CreateAccountPayload, UpdateAccountPayload } from '../types/accounts.ts'
import { accountQueryKeys } from './queryKeys.ts'

function useInvalidateAccountData() {
  const queryClient = useQueryClient()

  return () => {
    void queryClient.invalidateQueries({ queryKey: accountQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
  }
}

export function useCreateAccount() {
  const invalidateAccountData = useInvalidateAccountData()

  return useMutation({
    mutationFn: (payload: CreateAccountPayload) => createAccount(payload),
    onSuccess: invalidateAccountData,
  })
}

export function useUpdateAccount(accountId: string) {
  const invalidateAccountData = useInvalidateAccountData()

  return useMutation({
    mutationFn: (payload: UpdateAccountPayload) => updateAccount(accountId, payload),
    onSuccess: invalidateAccountData,
  })
}

export function useDeleteAccount() {
  const invalidateAccountData = useInvalidateAccountData()

  return useMutation({
    mutationFn: (accountId: string) => deleteAccount(accountId),
    onSuccess: invalidateAccountData,
  })
}
