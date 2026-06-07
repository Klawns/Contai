import type { RecurrenceFrequency } from './commitments.ts'

export type CommitmentFormValues = {
  description: string
  amount: number
  dueOn: string
  accountId: string
  categoryId: string
  note: string
  hasRecurrence: boolean
  recurrenceFrequency: RecurrenceFrequency
  recurrenceInterval: number
  recurrenceEndsOn: string
}
