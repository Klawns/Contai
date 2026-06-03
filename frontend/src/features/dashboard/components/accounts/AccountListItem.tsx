import { Landmark, Plus, WalletCards } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import type { AccountBalance } from '../../types/dashboard.ts'
import { getBalanceTone } from '../../utils/balance.ts'
import { formatCurrencyOrHidden } from '../../utils/formatters.ts'

type AccountListItemProps = {
  account: AccountBalance
  isBalanceHidden?: boolean
}

export function AccountListItem({
  account,
  isBalanceHidden = false,
}: AccountListItemProps) {
  const shouldReduceMotion = useReducedMotion()
  const tone = getBalanceTone(account.balance)
  const AccountIcon = account.type === 'wallet' ? WalletCards : Landmark

  return (
    <motion.article
      className="grid min-w-0 grid-cols-[38px_minmax(0,1fr)_36px] items-center gap-3 px-4 py-3.5"
      whileHover={shouldReduceMotion ? undefined : { x: 2 }}
      transition={{ duration: 0.16, ease: 'easeOut' }}
    >
      <span className="grid h-10 w-10 place-items-center rounded-full bg-[#f2eff8] text-[#6a22e5]">
        <AccountIcon className="h-5 w-5" aria-hidden="true" />
      </span>
      <div className="min-w-0">
        <strong className="block truncate text-[14px] font-semibold leading-tight text-[#241a30]">
          {account.name}
        </strong>
        <span
          className={`block truncate text-[12px] font-medium leading-tight ${
            isBalanceHidden ? 'text-[#b8b1c1]' : tone.textClass
          }`}
        >
          {formatCurrencyOrHidden(account.balance, isBalanceHidden)}
        </span>
      </div>
      <button
        type="button"
        className="grid h-9 w-9 cursor-pointer place-items-center rounded-full border-0 bg-[#6a22e5] text-white shadow-[0_8px_18px_rgba(104,24,232,0.24)] transition-transform hover:-translate-y-px focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        aria-label={`Adicionar lancamento em ${account.name}`}
      >
        <Plus className="h-5 w-5" aria-hidden="true" />
      </button>
    </motion.article>
  )
}
