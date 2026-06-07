import { useMemo } from 'react'
import { Plus } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { mapById } from '../../../../lib/collections/mapById.ts'
import { useActiveAccounts } from '../../../transactions/hooks/useActiveAccounts.ts'
import { useCreditCards } from '../../hooks/useCreditCards.ts'
import { PageShell } from '../shared/PageShell.tsx'
import { StatePanel } from '../shared/PageState.tsx'
import { CreditCardRow } from './CreditCardRow.tsx'

export function CreditCardListPage() {
  const navigate = useNavigate()
  const cardsQuery = useCreditCards()
  const accountsQuery = useActiveAccounts()
  const accountNames = useMemo(() => mapById(accountsQuery.data), [accountsQuery.data])

  return (
    <PageShell
      title="Cartoes"
      action={(
        <button
          type="button"
          className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22"
          aria-label="Novo cartao"
          onClick={() => navigate('/credit-cards/new')}
        >
          <Plus className="h-5 w-5" aria-hidden="true" />
        </button>
      )}
    >
      {cardsQuery.isLoading ? <StatePanel>Carregando cartoes...</StatePanel> : null}
      {cardsQuery.isError ? <StatePanel tone="danger">Nao foi possivel carregar os cartoes.</StatePanel> : null}
      {cardsQuery.data?.length === 0 ? <StatePanel>Ainda nao ha cartoes cadastrados.</StatePanel> : null}
      {cardsQuery.data?.length ? (
        <section className="bg-white">
          <ul className="divide-y divide-[#f0ebf6]">
            {cardsQuery.data.map((card) => (
              <CreditCardRow key={card.id} card={card} accounts={accountNames} />
            ))}
          </ul>
        </section>
      ) : null}
    </PageShell>
  )
}
