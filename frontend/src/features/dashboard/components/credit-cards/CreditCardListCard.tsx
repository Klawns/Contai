import { Link } from 'react-router-dom'
import type { AccountBalance, CreditCardDashboard } from '../../types/dashboard.ts'
import { DashboardEmptyStateCard } from '../cards/DashboardEmptyStateCard.tsx'
import { CreditCardListItem } from './CreditCardListItem.tsx'

type CreditCardListCardProps = {
  cards: CreditCardDashboard[]
  accounts: AccountBalance[]
  isBalanceHidden?: boolean
}

export function CreditCardListCard({
  cards,
  accounts,
  isBalanceHidden = false,
}: CreditCardListCardProps) {
  const accountById = new Map(accounts.map((account) => [account.accountId, account]))

  if (cards.length === 0) {
    return (
      <DashboardEmptyStateCard
        title="Ainda nao ha nenhum cartao"
        message="Cadastre cartoes para acompanhar limite e faturas sem alterar o saldo real."
        className="min-h-[190px]"
        action={(
          <Link
            to="/credit-cards/new"
            className="inline-flex min-h-10 items-center rounded-full bg-[#8f57ff] px-10 text-[14px] font-semibold text-white shadow-[0_10px_22px_rgba(143,87,255,0.24)] transition hover:bg-[#7e48ec] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#8f57ff]"
          >
            Adicionar cartao
          </Link>
        )}
      />
    )
  }

  return (
    <div className="overflow-hidden rounded-[18px] border border-[#ece8f2] bg-white shadow-[0_16px_38px_rgba(48,39,61,0.07)]">
      <div className="divide-y divide-[#f0edf5]">
        {cards.map((card) => (
          <CreditCardListItem
            key={card.cardId}
            card={card}
            account={accountById.get(card.linkedAccountId)}
            isBalanceHidden={isBalanceHidden}
          />
        ))}
      </div>
      <footer className="flex justify-center border-t border-[#ebe7f1] bg-[#fbfafe] px-4 py-4">
        <Link
          to="/credit-cards"
          className="inline-flex min-h-10 items-center rounded-full bg-[#8f57ff] px-18 text-[14px] font-semibold text-white shadow-[0_10px_22px_rgba(143,87,255,0.24)] transition hover:bg-[#7e48ec] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#8f57ff]"
        >
          Ver cartoes
        </Link>
      </footer>
    </div>
  )
}
