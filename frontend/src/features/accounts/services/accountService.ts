import { api } from '../../../lib/api/axios.ts'
import {
  accountsSchema,
  createAccountPayloadSchema,
  totalBalanceSchema,
  updateAccountPayloadSchema,
} from '../schemas/accounts.ts'
import type {
  Account,
  CreateAccountPayload,
  TotalBalance,
  UpdateAccountPayload,
} from '../types/accounts.ts'

export async function listActiveAccounts(): Promise<Account[]> {
  const response = await api.get<unknown>('/accounts', {
    params: { status: 'active' },
  })

  return accountsSchema.parse(response.data)
}

export async function getTotalBalance(): Promise<TotalBalance> {
  const response = await api.get<unknown>('/accounts/total-balance')

  return totalBalanceSchema.parse(response.data)
}

export async function createAccount(payload: CreateAccountPayload): Promise<Account> {
  const parsedPayload = createAccountPayloadSchema.parse(payload)
  const response = await api.post<unknown>('/accounts', parsedPayload)

  return accountsSchema.element.parse(response.data)
}

export async function updateAccount(
  accountId: string,
  payload: UpdateAccountPayload,
): Promise<Account> {
  const parsedPayload = updateAccountPayloadSchema.parse(payload)
  const response = await api.patch<unknown>(`/accounts/${accountId}`, parsedPayload)

  return accountsSchema.element.parse(response.data)
}

export async function deleteAccount(accountId: string): Promise<void> {
  await api.delete(`/accounts/${accountId}`)
}
