import type {
  CategoryTransactionType,
  TransactionFilters,
} from '../types/transactions.ts'

export const transactionsQueryKeys = {
  all: ['transactions'] as const,
  list: (filters: TransactionFilters) =>
    [
      ...transactionsQueryKeys.all,
      'list',
      filters.startAt ?? null,
      filters.endAt ?? null,
      filters.accountId ?? null,
      filters.categoryId ?? null,
      filters.type ?? null,
      filters.limit ?? null,
      filters.offset ?? null,
    ] as const,
}

export const accountsQueryKeys = {
  all: ['accounts'] as const,
}

export const categoriesQueryKeys = {
  all: ['categories'] as const,
  list: (type: CategoryTransactionType) =>
    [...categoriesQueryKeys.all, 'list', type] as const,
}
