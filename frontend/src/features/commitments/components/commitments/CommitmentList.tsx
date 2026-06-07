import { useNavigate } from 'react-router-dom'
import { useConfirmDialog } from '../../../../components/confirm-dialog-context.ts'
import type { Account, Category } from '../../../transactions/types/transactions.ts'
import { useCancelCommitment } from '../../hooks/useCommitmentMutations.ts'
import type { Commitment, CommitmentType } from '../../types/commitments.ts'
import { CommitmentListItem } from './CommitmentListItem.tsx'
import { EmptyCommitmentState } from './EmptyCommitmentState.tsx'

export function CommitmentList({
  type,
  commitments,
  accountNames,
  categoryNames,
}: {
  type: CommitmentType
  commitments: Commitment[]
  accountNames: Map<string, Account>
  categoryNames: Map<string, Category>
}) {
  const navigate = useNavigate()
  const { confirm } = useConfirmDialog()
  const cancelMutation = useCancelCommitment(type)

  async function handleCancel(commitment: Commitment) {
    const shouldCancel = await confirm({
      title: 'Cancelar compromisso',
      description: `Cancelar "${commitment.description}"?`,
      confirmLabel: 'Cancelar compromisso',
      cancelLabel: 'Voltar',
      tone: 'danger',
    })

    if (shouldCancel) {
      cancelMutation.mutate(commitment.id)
    }
  }

  if (commitments.length === 0) {
    return <EmptyCommitmentState type={type} />
  }

  return (
    <section className="bg-white">
      <ul className="divide-y divide-[#f0ebf6]">
        {commitments.map((commitment) => (
          <CommitmentListItem
            key={commitment.id}
            type={type}
            commitment={commitment}
            account={accountNames.get(commitment.accountId)}
            category={categoryNames.get(commitment.categoryId)}
            isCancelling={cancelMutation.isPending}
            onCancel={(item) => void handleCancel(item)}
            onEdit={(item) =>
              navigate(`/planning/edit?type=${type}&id=${encodeURIComponent(item.id)}`)
            }
            onSettle={(item) =>
              navigate(`/planning/settle?type=${type}&id=${encodeURIComponent(item.id)}`)
            }
          />
        ))}
      </ul>
    </section>
  )
}
