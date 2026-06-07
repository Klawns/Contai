import type { CategoryTransactionType } from '../../transactions/types/transactions.ts'
import type { CommitmentType } from '../types/commitments.ts'

export function parseCommitmentType(value: string | null): CommitmentType {
  return value === 'receivable' ? 'receivable' : 'payable'
}

export function categoryTypeForCommitment(type: CommitmentType): CategoryTransactionType {
  return type === 'payable' ? 'expense' : 'income'
}
