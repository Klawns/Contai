import { Link } from 'react-router-dom'
import { CreditCard } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import type { AccountBalance, CreditCardDashboard } from '../../types/dashboard.ts'
import { formatCurrencyOrHidden } from '../../utils/formatters.ts'

type CreditCardListItemProps = {
  card: CreditCardDashboard
  account?: AccountBalance
  isBalanceHidden?: boolean
}

function formatDate(value: string | null) {
  if (!value) {
    return 'Sem vencimento'
  }

  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'short',
  }).format(new Date(value))
}

export function CreditCardListItem({
  card,
  account,
  isBalanceHidden = false,
}: CreditCardListItemProps) {
  const shouldReduceMotion = useReducedMotion()

  return (
    <motion.article
      className="grid min-w-0 grid-cols-[38px_minmax(0,1fr)_minmax(92px,auto)] items-center gap-3 px-4 py-3.5"
      whileHover={shouldReduceMotion ? undefined : { x: 2 }}
      transition={{ duration: 0.16, ease: 'easeOut' }}
    >
      <span className="grid h-10 w-10 place-items-center rounded-full bg-[#eef6ff] text-[#216fb8]">
        <CreditCard className="h-5 w-5" aria-hidden="true" />
      </span>
      <div className="min-w-0">
        <strong className="block truncate text-[14px] font-semibold leading-tight text-[#241a30]">
          {card.name}
        </strong>
        <span className="block truncate text-[12px] font-medium leading-tight text-[#81788c]">
          {account?.name ?? 'Conta vinculada'} / livre {formatCurrencyOrHidden(card.limitAvailable, isBalanceHidden)}
        </span>
      </div>
      <Link
        to={card.currentInvoiceId ? `/credit-cards/invoice?invoiceId=${encodeURIComponent(card.currentInvoiceId)}` : '/credit-cards'}
        className="min-w-0 text-right focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
      >
        <strong className={`block truncate text-[13px] font-semibold leading-tight ${isBalanceHidden ? 'text-[#b8b1c1]' : 'text-[#c72f4d]'}`}>
          {formatCurrencyOrHidden(card.currentInvoiceAmount, isBalanceHidden)}
        </strong>
        <span className="mt-1 block truncate text-[12px] font-semibold leading-tight text-[#958c9f]">
          {formatDate(card.currentInvoiceDueAt)}
        </span>
      </Link>
    </motion.article>
  )
}
