import { toDateInputValue } from '../../transactions/utils/date.ts'
import type { Commitment, SettlementPayload } from '../types/commitments.ts'
import type { SettlementFormValues } from '../types/settlementForm.ts'
import { toLocalDateTime } from './commitmentDates.ts'

export function getSettlementDefaultValues(commitment: Commitment): SettlementFormValues {
  return {
    amount: commitment.amount,
    occurredOn: toDateInputValue(new Date()),
    accountId: commitment.accountId,
    categoryId: commitment.categoryId,
    note: commitment.note,
  }
}

export function toSettlementPayload(values: SettlementFormValues): SettlementPayload {
  return {
    amount: values.amount,
    occurredAt: toLocalDateTime(values.occurredOn),
    accountId: values.accountId,
    categoryId: values.categoryId,
    note: values.note.trim(),
  }
}
