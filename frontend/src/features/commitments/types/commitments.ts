export const commitmentTypes = ['payable', 'receivable'] as const
export const commitmentStatuses = ['pending', 'settled', 'canceled'] as const
export const effectiveCommitmentStatuses = [
  'pending',
  'overdue',
  'paid',
  'received',
  'canceled',
] as const
export const recurrenceFrequencies = ['daily', 'weekly', 'monthly', 'yearly'] as const

export type CommitmentType = (typeof commitmentTypes)[number]
export type CommitmentStatus = (typeof commitmentStatuses)[number]
export type EffectiveCommitmentStatus = (typeof effectiveCommitmentStatuses)[number]
export type RecurrenceFrequency = (typeof recurrenceFrequencies)[number]

export type CommitmentRecurrence = {
  frequency: RecurrenceFrequency
  interval: number
  endsAt: string | null
}

export type Commitment = {
  id: string
  userId: string
  type: CommitmentType
  description: string
  amount: number
  dueAt: string
  accountId: string
  categoryId: string
  note: string
  status: CommitmentStatus
  effectiveStatus: EffectiveCommitmentStatus
  recurrence: CommitmentRecurrence | null
  settledAt: string | null
  settlementTransactionId: string | null
  canceledAt: string | null
  createdAt: string
  updatedAt: string
}

export type CommitmentFilters = {
  startAt?: string
  endAt?: string
  status?: CommitmentStatus
  effectiveStatus?: EffectiveCommitmentStatus
  accountId?: string
  categoryId?: string
  limit?: number
  offset?: number
}

export type CommitmentPayload = {
  description: string
  amount: number
  dueAt: string
  accountId: string
  categoryId: string
  note: string
  recurrence: CommitmentRecurrence | null
}

export type SettlementPayload = {
  amount: number
  occurredAt: string
  accountId: string
  categoryId: string
  note: string
}
