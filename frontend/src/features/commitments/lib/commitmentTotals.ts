import type { Commitment } from '../types/commitments.ts'

export function calculateCommitmentTotals(commitments: Commitment[] | undefined) {
  return (commitments ?? []).reduce(
    (accumulator, commitment) => {
      if (commitment.effectiveStatus === 'pending' || commitment.effectiveStatus === 'overdue') {
        accumulator.open += commitment.amount
      }
      if (commitment.effectiveStatus === 'overdue') {
        accumulator.overdue += commitment.amount
      }
      return accumulator
    },
    { open: 0, overdue: 0 },
  )
}
