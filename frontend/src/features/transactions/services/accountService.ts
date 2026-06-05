import { api } from '../../../lib/api/axios.ts'
import { accountsSchema } from '../schemas/transactions.ts'
import type { Account } from '../types/transactions.ts'

export async function listActiveAccounts(): Promise<Account[]> {
  const response = await api.get<unknown>('/accounts', {
    params: { status: 'active' },
  })

  return accountsSchema.parse(response.data)
}
