import type { AccountBalance } from '../../types/dashboard.ts'
import { DashboardEmptyStateCard } from '../DashboardEmptyStateCard.tsx'
import { AccountListItem } from './AccountListItem.tsx'

type AccountListCardProps = {
  accounts: AccountBalance[]
  isBalanceHidden?: boolean
}

export function AccountListCard({
  accounts,
  isBalanceHidden = false,
}: AccountListCardProps) {
  if (accounts.length === 0) {
    return (
      <DashboardEmptyStateCard
        title="Ainda nao ha nenhuma conta"
        message="Crie uma conta para acompanhar seus saldos no dashboard."
        className="min-h-[214px]"
        action={
          <button
            type="button"
            className="min-h-10 rounded-full bg-[#8f57ff] px-10 text-[14px] font-semibold text-white shadow-[0_10px_22px_rgba(143,87,255,0.24)] transition hover:bg-[#7e48ec] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#8f57ff]"
          >
            Adicionar conta
          </button>
        }
      />
    )
  }

  return (
    <div className="overflow-hidden rounded-[18px] border border-[#ece8f2] bg-white shadow-[0_16px_38px_rgba(48,39,61,0.07)]">
      <div className="divide-y divide-[#f0edf5]">
        {accounts.map((account) => (
          <AccountListItem
            key={account.accountId}
            account={account}
            isBalanceHidden={isBalanceHidden}
          />
        ))}
      </div>
      <footer className="flex justify-center border-t border-[#ebe7f1] bg-[#fbfafe] px-4 py-4">
        <button
          type="button"
          className="min-h-10 rounded-full bg-[#8f57ff] px-18 text-[14px] font-semibold text-white shadow-[0_10px_22px_rgba(143,87,255,0.24)] transition hover:bg-[#7e48ec] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#8f57ff]"
        >
          Adicionar conta
        </button>
      </footer>
    </div>
  )
}
