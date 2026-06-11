export const transactionTypes = ['income', 'expense', 'transfer'] as const
export const categoryTransactionTypes = ['income', 'expense'] as const
export const transactionOriginTypes = ['manual', 'payable', 'receivable', 'credit_card_invoice'] as const
export const settlementStatuses = ['settled', 'pending'] as const
export const recurrenceTypes = ['none', 'fixed', 'repeat'] as const
export const recurrenceFrequencies = ['daily', 'weekly', 'monthly'] as const

export type TransactionType = (typeof transactionTypes)[number]
export type CategoryTransactionType = (typeof categoryTransactionTypes)[number]
export type TransactionOriginType = (typeof transactionOriginTypes)[number]
export type SettlementStatus = (typeof settlementStatuses)[number]
export type RecurrenceType = (typeof recurrenceTypes)[number]
export type RecurrenceFrequency = (typeof recurrenceFrequencies)[number]

export type TransactionRecurrence = {
  frequency: RecurrenceFrequency
  quantity: number | null
  startsAt: string
  endsAt: string | null
  dayOfMonth: number | null
}

export type Transaction = {
  id: string
  userId: string
  type: TransactionType
  description: string
  amount: number
  occurredAt: string
  accountId: string | null
  sourceAccountId: string | null
  destinationAccountId: string | null
  categoryId: string | null
  status: string
  originType: TransactionOriginType
  originId: string | null
  settlementStatus: SettlementStatus
  settledAt: string | null
  recurrenceType: RecurrenceType
  recurrence: TransactionRecurrence | null
  note: string
  removedAt: string | null
  createdAt: string
  updatedAt: string
}

export type Account = {
  id: string
  userId: string
  name: string
  type: string
  initialBalance: number
  currentBalance: number
  bankIconId: string
  includeInDashboardTotal: boolean
  status: string
  createdAt: string
  updatedAt: string
}

export type Category = {
  id: string
  userId: string
  name: string
  normalizedName: string
  type: CategoryTransactionType
  color: string
  icon: string
  isDefault: boolean
  status: string
  createdAt: string
  updatedAt: string
}

export type TransactionFilters = {
  startAt?: string
  endAt?: string
  accountId?: string
  categoryId?: string
  type?: TransactionType
  settlementStatus?: SettlementStatus
  limit?: number
  offset?: number
}

export type CreateIncomeTransactionPayload = {
  description: string
  amount: number
  occurredAt: string
  accountId: string | null
  categoryId: string
  settlementStatus: SettlementStatus
  settledAt: string | null
  recurrenceType: RecurrenceType
  recurrence: TransactionRecurrence | null
  note: string
}

export type CreateExpenseTransactionPayload = CreateIncomeTransactionPayload

export type CreateTransferTransactionPayload = {
  description: string
  amount: number
  occurredAt: string
  sourceAccountId: string
  destinationAccountId: string
  note: string
}

export type CreateTransactionPayloadByType = {
  income: CreateIncomeTransactionPayload
  expense: CreateExpenseTransactionPayload
  transfer: CreateTransferTransactionPayload
}

export type UpdateIncomeExpenseTransactionPayload = CreateIncomeTransactionPayload

export type UpdateTransferTransactionPayload = CreateTransferTransactionPayload

export type UpdateTransactionPayloadByType = {
  income: UpdateIncomeExpenseTransactionPayload
  expense: UpdateIncomeExpenseTransactionPayload
  transfer: UpdateTransferTransactionPayload
}

export type UpdateTransactionPayload =
  UpdateTransactionPayloadByType[keyof UpdateTransactionPayloadByType]

export type CreateCategoryPayload = {
  name: string
  color: string
  icon: string
  type: CategoryTransactionType
}
