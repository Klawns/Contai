import type { CommitmentFilters, CommitmentType } from '../types/commitments.ts'

export const commitmentsQueryKeys = {
  all: ['commitments'] as const,
  list: (type: CommitmentType, filters: CommitmentFilters) =>
    [
      ...commitmentsQueryKeys.all,
      type,
      filters.startAt ?? null,
      filters.endAt ?? null,
      filters.status ?? null,
      filters.effectiveStatus ?? null,
      filters.accountId ?? null,
      filters.categoryId ?? null,
      filters.limit ?? null,
      filters.offset ?? null,
    ] as const,
}
