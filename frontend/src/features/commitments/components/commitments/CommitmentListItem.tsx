import { Clock, Landmark, PencilLine, XCircle } from 'lucide-react'
import type { Account, Category } from '../../../transactions/types/transactions.ts'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { formatCommitmentDate } from '../../lib/commitmentDates.ts'
import { typeCopy } from '../../lib/commitmentPresentation.ts'
import type { Commitment, CommitmentType } from '../../types/commitments.ts'
import { CommitmentBadge } from './CommitmentBadge.tsx'

type CommitmentListItemProps = {
  type: CommitmentType
  commitment: Commitment
  account?: Account
  category?: Category
  isCancelling: boolean
  onCancel: (commitment: Commitment) => void
  onEdit: (commitment: Commitment) => void
  onSettle: (commitment: Commitment) => void
}

export function CommitmentListItem({
  type,
  commitment,
  account,
  category,
  isCancelling,
  onCancel,
  onEdit,
  onSettle,
}: CommitmentListItemProps) {
  const TypeIcon = typeCopy[type].icon
  const isPending =
    commitment.effectiveStatus === 'pending' || commitment.effectiveStatus === 'overdue'

  return (
    <li className="grid gap-3 px-1 py-3 sm:px-2">
      <div className="grid grid-cols-[40px_minmax(0,1fr)_auto] items-center gap-3">
        <span
          className={`grid h-10 w-10 flex-none place-items-center rounded-full ${
            type === 'payable'
              ? 'bg-[#fff0f2] text-[#c72f4d]'
              : 'bg-[#e8f8ef] text-[#147a46]'
          }`}
        >
          <TypeIcon className="h-4.5 w-4.5" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h3 className="truncate text-[14px] font-semibold leading-tight text-[#2c2237]">
            {commitment.description}
          </h3>
          <p className="mt-1 flex min-w-0 items-center gap-1 text-[12px] font-semibold leading-tight text-[#81788c]">
            <Clock className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
            <span className="truncate">{formatCommitmentDate(commitment.dueAt)}</span>
          </p>
        </div>
        <div className="min-w-0 text-right">
          <strong
            className={`block text-[14px] font-semibold leading-tight ${
              type === 'payable' ? 'text-[#c72f4d]' : 'text-[#147a46]'
            }`}
          >
            {formatCurrency(commitment.amount)}
          </strong>
          <span className="mt-1 block">
            <CommitmentBadge commitment={commitment} />
          </span>
        </div>
      </div>

      <div className="flex min-w-0 items-center gap-2 pl-[52px] text-[12px] font-semibold text-[#81788c]">
        <Landmark className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
        <span className="truncate">{[account?.name, category?.name].filter(Boolean).join(' / ')}</span>
      </div>

      {isPending ? (
        <div className="flex flex-wrap gap-2 pl-[52px]">
          <button
            type="button"
            className="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-full bg-[#281d35] px-3 text-[12px] font-semibold text-white transition-colors hover:bg-[#3a2a4a] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
            onClick={() => onSettle(commitment)}
          >
            <TypeIcon className="h-3.5 w-3.5" aria-hidden="true" />
            {typeCopy[type].settleButton}
          </button>
          <button
            type="button"
            className="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-full border border-[#e3ddea] bg-white px-3 text-[12px] font-semibold text-[#4f435c] transition-colors hover:bg-[#f8f5fb] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
            onClick={() => onEdit(commitment)}
          >
            <PencilLine className="h-3.5 w-3.5" aria-hidden="true" />
            Editar
          </button>
          <button
            type="button"
            disabled={isCancelling}
            className="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-full border border-[#f0caca] bg-white px-3 text-[12px] font-semibold text-[#c75959] transition-colors hover:bg-[#fff4f4] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#c75959] disabled:cursor-not-allowed disabled:opacity-55"
            onClick={() => onCancel(commitment)}
          >
            <XCircle className="h-3.5 w-3.5" aria-hidden="true" />
            Cancelar
          </button>
        </div>
      ) : null}
    </li>
  )
}
