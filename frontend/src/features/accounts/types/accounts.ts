export const accountTypes = [
  'checking',
  'savings',
  'digital',
  'cash',
  'salary',
  'investment',
  'other',
] as const

export type AccountType = (typeof accountTypes)[number]

export type Account = {
  id: string
  userId: string
  name: string
  type: AccountType
  initialBalance: number
  currentBalance: number
  bankIconId: string
  includeInDashboardTotal: boolean
  status: string
  createdAt: string
  updatedAt: string
}

export type TotalBalance = {
  totalBalance: number
}

export type CreateAccountPayload = {
  name: string
  type: AccountType
  initialBalance: number
  bankIconId: string
  includeInDashboardTotal: boolean
}

export type UpdateAccountPayload = Omit<CreateAccountPayload, 'initialBalance'>
