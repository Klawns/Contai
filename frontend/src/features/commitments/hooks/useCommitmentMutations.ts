import { useMutation, useQueryClient } from '@tanstack/react-query'
import { accountQueryKeys } from '../../accounts/hooks/queryKeys.ts'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { accountsQueryKeys, transactionsQueryKeys } from '../../transactions/hooks/queryKeys.ts'
import {
  cancelCommitment,
  createCommitment,
  settleCommitment,
  updateCommitment,
} from '../services/commitmentService.ts'
import type {
  CommitmentPayload,
  CommitmentType,
  SettlementPayload,
} from '../types/commitments.ts'
import { commitmentsQueryKeys } from './queryKeys.ts'

function useInvalidateFinancialData() {
  const queryClient = useQueryClient()

  return () => {
    void queryClient.invalidateQueries({ queryKey: commitmentsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: transactionsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: dashboardQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: accountsQueryKeys.all })
    void queryClient.invalidateQueries({ queryKey: accountQueryKeys.all })
  }
}

export function useCreateCommitment(type: CommitmentType) {
  const invalidate = useInvalidateFinancialData()

  return useMutation({
    mutationFn: (payload: CommitmentPayload) => createCommitment(type, payload),
    onSuccess: invalidate,
  })
}

export function useUpdateCommitment(type: CommitmentType, commitmentId: string) {
  const invalidate = useInvalidateFinancialData()

  return useMutation({
    mutationFn: (payload: CommitmentPayload) =>
      updateCommitment(type, commitmentId, payload),
    onSuccess: invalidate,
  })
}

export function useSettleCommitment(type: CommitmentType, commitmentId: string) {
  const invalidate = useInvalidateFinancialData()

  return useMutation({
    mutationFn: (payload: SettlementPayload) => settleCommitment(type, commitmentId, payload),
    onSuccess: invalidate,
  })
}

export function useCancelCommitment(type: CommitmentType) {
  const invalidate = useInvalidateFinancialData()

  return useMutation({
    mutationFn: (commitmentId: string) => cancelCommitment(type, commitmentId),
    onSuccess: invalidate,
  })
}
