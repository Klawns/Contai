import { api } from '../../../lib/api/axios.ts'
import {
  createIncomeExpenseTransactionPayloadSchema,
  createTransferTransactionPayloadSchema,
  transactionFiltersSchema,
  transactionSchema,
  transactionsSchema,
  updateIncomeExpenseTransactionPayloadSchema,
  updateTransferTransactionPayloadSchema,
} from '../schemas/transactions.ts'
import type {
  CreateExpenseTransactionPayload,
  CreateIncomeTransactionPayload,
  CreateTransferTransactionPayload,
  Transaction,
  TransactionFilters,
  TransactionType,
  UpdateTransactionPayloadByType,
} from '../types/transactions.ts'

export async function listTransactions(
  filters: TransactionFilters,
): Promise<Transaction[]> {
  const params = transactionFiltersSchema.parse(filters)
  const response = await api.get<unknown>('/transactions', { params })

  return transactionsSchema.parse(response.data)
}

export async function createIncomeTransaction(
  payload: CreateIncomeTransactionPayload,
): Promise<Transaction> {
  const body = createIncomeExpenseTransactionPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/transactions/income', body)

  return transactionSchema.parse(response.data)
}

export async function createExpenseTransaction(
  payload: CreateExpenseTransactionPayload,
): Promise<Transaction> {
  const body = createIncomeExpenseTransactionPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/transactions/expense', body)

  return transactionSchema.parse(response.data)
}

export async function createTransferTransaction(
  payload: CreateTransferTransactionPayload,
): Promise<Transaction> {
  const body = createTransferTransactionPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/transactions/transfer', body)

  return transactionSchema.parse(response.data)
}

export async function updateTransaction<TType extends TransactionType>(
  type: TType,
  transactionId: string,
  payload: UpdateTransactionPayloadByType[TType],
): Promise<Transaction> {
  const body =
    type === 'transfer'
      ? updateTransferTransactionPayloadSchema.parse(payload)
      : updateIncomeExpenseTransactionPayloadSchema.parse(payload)
  const response = await api.patch<unknown>(`/transactions/${transactionId}`, body)

  return transactionSchema.parse(response.data)
}

export async function deleteTransaction(transactionId: string): Promise<void> {
  await api.delete(`/transactions/${transactionId}`)
}
