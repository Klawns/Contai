import { ArrowDownLeft, ArrowUpRight, Repeat2 } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import { getBalanceTone } from '../../utils/balance.ts'
import { formatCurrencyOrHidden } from '../../utils/formatters.ts'

type MonthlyBalanceCardProps = {
  monthlyIncome: number
  monthlyExpense: number
  monthlyNetBalance: number
  monthlyTransferIn: number
  monthlyTransferOut: number
  isBalanceHidden?: boolean
}

export function MonthlyBalanceCard({
  monthlyIncome,
  monthlyExpense,
  monthlyNetBalance,
  monthlyTransferIn,
  monthlyTransferOut,
  isBalanceHidden = false,
}: MonthlyBalanceCardProps) {
  const shouldReduceMotion = useReducedMotion()
  const netTone = getBalanceTone(monthlyNetBalance)
  const transferTotal = monthlyTransferIn + monthlyTransferOut
  const rows = [
    {
      label: 'Receitas',
      value: monthlyIncome,
      className: 'text-[#159b58]',
      icon: ArrowUpRight,
      iconClassName: 'bg-[#e9f8ef] text-[#17a760]',
    },
    {
      label: 'Despesas',
      value: monthlyExpense,
      className: 'text-[#d83b3b]',
      icon: ArrowDownLeft,
      iconClassName: 'bg-[#fdecec] text-[#e44545]',
    },
    {
      label: 'Transferencias',
      value: transferTotal,
      className: 'text-[#6a22e5]',
      icon: Repeat2,
      iconClassName: 'bg-[#f2eff8] text-[#6a22e5]',
    },
  ]

  return (
    <motion.article
      className="min-w-[calc(100vw-32px)] snap-center rounded-[18px] border border-[#ece8f2] bg-white p-4 shadow-[0_16px_38px_rgba(48,39,61,0.07)] sm:min-w-[360px] md:min-w-0"
      whileHover={shouldReduceMotion ? undefined : { y: -2 }}
      whileTap={shouldReduceMotion ? undefined : { scale: 0.99 }}
      transition={{ duration: 0.18, ease: 'easeOut' }}
    >
      <h3 className="m-0 text-[15px] font-semibold leading-tight text-[#241a30]">
        Balanco mensal
      </h3>
      <div className="mt-4 grid divide-y divide-[#f0edf5]">
        {rows.map((row) => {
          const Icon = row.icon

          return (
            <div key={row.label} className="grid grid-cols-[34px_minmax(0,1fr)_auto] items-center gap-2.5 py-3 first:pt-0">
              <span
                className={`grid h-8 w-8 place-items-center rounded-full ${row.iconClassName}`}
                aria-hidden="true"
              >
                <Icon className="h-4 w-4" />
              </span>
              <span className="truncate text-[13px] font-medium text-[#6d6478]">
                {row.label}
              </span>
              <strong
                className={`truncate text-[13px] font-semibold ${
                  isBalanceHidden ? 'text-[#b8b1c1]' : row.className
                }`}
              >
                {formatCurrencyOrHidden(row.value, isBalanceHidden)}
              </strong>
            </div>
          )
        })}
      </div>
      <div className="mt-3 rounded-xl bg-[#fbfafe] px-3 py-3">
        <span className="block text-[12px] font-medium text-[#8b8394]">
          Balanco do mes
        </span>
        <strong
          className={`block truncate text-[24px] font-semibold leading-tight ${
            isBalanceHidden ? 'text-[#b8b1c1]' : netTone.textClass
          }`}
        >
          {formatCurrencyOrHidden(monthlyNetBalance, isBalanceHidden)}
        </strong>
      </div>
    </motion.article>
  )
}
