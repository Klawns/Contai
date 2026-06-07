import { statusCopy } from '../../lib/commitmentPresentation.ts'
import type { Commitment } from '../../types/commitments.ts'

export function CommitmentBadge({ commitment }: { commitment: Commitment }) {
  const tone =
    commitment.effectiveStatus === 'overdue'
      ? 'bg-[#fff0f2] text-[#c72f4d]'
      : commitment.effectiveStatus === 'paid' || commitment.effectiveStatus === 'received'
        ? 'bg-[#e8f8ef] text-[#147a46]'
        : commitment.effectiveStatus === 'canceled'
          ? 'bg-[#f1eef5] text-[#81788c]'
          : 'bg-[#fff8e8] text-[#9b6a12]'

  return (
    <span className={`rounded-full px-2.5 py-1 text-[11px] font-semibold ${tone}`}>
      {statusCopy[commitment.effectiveStatus]}
    </span>
  )
}
