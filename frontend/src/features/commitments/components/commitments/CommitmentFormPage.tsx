import { useMemo } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useCommitments } from '../../hooks/useCommitments.ts'
import { typeCopy } from '../../lib/commitmentPresentation.ts'
import type { CommitmentType } from '../../types/commitments.ts'
import { CommitmentPageState } from '../shared/CommitmentPageState.tsx'
import { CommitmentForm } from './CommitmentForm.tsx'

export function CommitmentFormPage({
  type,
  mode = 'create',
}: {
  type: CommitmentType
  mode?: 'create' | 'edit'
}) {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const commitmentId = searchParams.get('id') ?? ''
  const commitmentsQuery = useCommitments(type, {})
  const commitment = useMemo(
    () =>
      mode === 'edit' && commitmentId
        ? commitmentsQuery.data?.find((item) => item.id === commitmentId)
        : undefined,
    [commitmentId, commitmentsQuery.data, mode],
  )

  if (mode === 'edit' && commitmentsQuery.isLoading) {
    return (
      <CommitmentPageState
        animationKey="commitment-edit-loading"
        tone={typeCopy[type].tone}
        message="Carregando compromisso..."
      />
    )
  }

  if (mode === 'edit' && (!commitment || commitment.status !== 'pending')) {
    return (
      <CommitmentPageState
        animationKey="commitment-edit-invalid"
        tone={typeCopy[type].tone}
        message="Este compromisso nao pode ser editado."
        danger
        onBack={() => navigate('/planning')}
      />
    )
  }

  return <CommitmentForm type={type} mode={mode} initialCommitment={commitment} />
}
