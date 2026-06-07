import { toDateInputValue } from '../../transactions/utils/date.ts'
import type {
  Commitment,
  CommitmentPayload,
  CommitmentRecurrence,
} from '../types/commitments.ts'
import type { CommitmentFormValues } from '../types/commitmentForm.ts'
import { toLocalDateTime } from './commitmentDates.ts'

export function getDefaultFormValues(): CommitmentFormValues {
  return {
    description: '',
    amount: 0,
    dueOn: toDateInputValue(new Date()),
    accountId: '',
    categoryId: '',
    note: '',
    hasRecurrence: false,
    recurrenceFrequency: 'monthly',
    recurrenceInterval: 1,
    recurrenceEndsOn: '',
  }
}

export function getInitialFormValues(commitment: Commitment): CommitmentFormValues {
  return {
    description: commitment.description,
    amount: commitment.amount,
    dueOn: toDateInputValue(new Date(commitment.dueAt)),
    accountId: commitment.accountId,
    categoryId: commitment.categoryId,
    note: commitment.note,
    hasRecurrence: Boolean(commitment.recurrence),
    recurrenceFrequency: commitment.recurrence?.frequency ?? 'monthly',
    recurrenceInterval: commitment.recurrence?.interval ?? 1,
    recurrenceEndsOn: commitment.recurrence?.endsAt
      ? toDateInputValue(new Date(commitment.recurrence.endsAt))
      : '',
  }
}

export function toCommitmentPayload(values: CommitmentFormValues): CommitmentPayload {
  const recurrence: CommitmentRecurrence | null = values.hasRecurrence
    ? {
        frequency: values.recurrenceFrequency,
        interval: values.recurrenceInterval,
        endsAt: values.recurrenceEndsOn ? toLocalDateTime(values.recurrenceEndsOn) : null,
      }
    : null

  return {
    description: values.description.trim(),
    amount: values.amount,
    dueAt: toLocalDateTime(values.dueOn),
    accountId: values.accountId,
    categoryId: values.categoryId,
    note: values.note.trim(),
    recurrence,
  }
}
