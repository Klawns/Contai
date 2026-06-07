import { useMemo } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useCommitments } from '../../hooks/useCommitments.ts'
import { typeCopy } from '../../lib/commitmentPresentation.ts'
import { parseCommitmentType } from '../../lib/commitmentType.ts'
import { CommitmentPageState } from '../shared/CommitmentPageState.tsx'
import { SettlementForm } from './SettlementForm.tsx'

export function SettlementPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const type = parseCommitmentType(searchParams.get('type'))
  const commitmentId = searchParams.get('id') ?? ''
  const commitmentsQuery = useCommitments(type, {})
  const commitment = useMemo(
    () => commitmentsQuery.data?.find((item) => item.id === commitmentId),
    [commitmentId, commitmentsQuery.data],
  )

  if (commitmentsQuery.isLoading) {
    return (
      <CommitmentPageState
        animationKey="settle-loading"
        tone={typeCopy[type].tone}
        message="Carregando compromisso..."
      />
    )
  }

  if (!commitment || commitment.status !== 'pending') {
    return (
      <CommitmentPageState
        animationKey="settle-invalid"
        tone={typeCopy[type].tone}
        message="Este compromisso nao pode ser quitado."
        danger
        onBack={() => navigate('/planning')}
      />
    )
  }

  return <SettlementForm type={type} commitment={commitment} />
}
