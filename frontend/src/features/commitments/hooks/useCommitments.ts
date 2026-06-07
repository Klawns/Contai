import { useQuery } from '@tanstack/react-query'
import { listCommitments } from '../services/commitmentService.ts'
import type { CommitmentFilters, CommitmentType } from '../types/commitments.ts'
import { commitmentsQueryKeys } from './queryKeys.ts'

export function useCommitments(type: CommitmentType, filters: CommitmentFilters) {
  return useQuery({
    queryKey: commitmentsQueryKeys.list(type, filters),
    queryFn: () => listCommitments(type, filters),
  })
}
