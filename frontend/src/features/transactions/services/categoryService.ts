import { api } from '../../../lib/api/axios.ts'
import {
  categoriesSchema,
  categorySchema,
  categoryTransactionTypeSchema,
  createCategoryPayloadSchema,
} from '../schemas/transactions.ts'
import type {
  Category,
  CategoryTransactionType,
  CreateCategoryPayload,
} from '../types/transactions.ts'

export async function listActiveCategories(
  type: CategoryTransactionType,
): Promise<Category[]> {
  const categoryType = categoryTransactionTypeSchema.parse(type)
  const response = await api.get<unknown>('/categories', {
    params: { type: categoryType, status: 'active' },
  })

  return categoriesSchema.parse(response.data)
}

export async function createCategory(payload: CreateCategoryPayload): Promise<Category> {
  const body = createCategoryPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/categories', body)

  return categorySchema.parse(response.data)
}
