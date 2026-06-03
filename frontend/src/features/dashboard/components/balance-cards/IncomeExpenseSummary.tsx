import type { LucideIcon } from 'lucide-react'
import { ArrowDown, ArrowUp } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'

type IncomeExpenseSummaryProps = {
  income: number
  expense: number
  isHidden?: boolean
}

type SummaryItem = {
  label: string
  value: number
  variant: 'income' | 'expense'
  icon: LucideIcon
}

const variantClasses = {
  income: {
    icon: 'bg-[#43c762] text-white',
    value: 'text-[#43c762]',
  },
  expense: {
    icon: 'bg-[#ff4646] text-white',
    value: 'text-[#ff4646]',
  },
}

const summaryValueFormatter = new Intl.NumberFormat('pt-BR', {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

function formatSummaryValue(valueInCents: number) {
  const absolute = Math.abs(valueInCents) / 100
  const formatted = summaryValueFormatter.format(absolute)

  return valueInCents < 0 ? `-${formatted}` : formatted
}

export function IncomeExpenseSummary({
  income,
  expense,
  isHidden = false,
}: IncomeExpenseSummaryProps) {
  const shouldReduceMotion = useReducedMotion()
  const valueTransition = { duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' as const }
  const items: SummaryItem[] = [
    { label: 'Receitas', value: income, variant: 'income', icon: ArrowUp },
    { label: 'Despesas', value: expense, variant: 'expense', icon: ArrowDown },
  ]

  return (
    <div className="mx-auto grid w-full max-w-[340px] grid-cols-2 gap-3 sm:max-w-[380px] sm:gap-5">
      {items.map((item, index) => {
        const Icon = item.icon
        const classes = variantClasses[item.variant]
        const alignmentClass = index === 0 ? 'justify-self-start' : 'justify-self-end'

        return (
          <div
            key={item.label}
            className={`grid w-fit max-w-full min-w-0 grid-cols-[48px_minmax(0,1fr)] items-center gap-3 sm:grid-cols-[52px_minmax(0,1fr)] ${alignmentClass}`}
          >
            <span
              className={`grid h-12 w-12 place-items-center rounded-full sm:h-[52px] sm:w-[52px] ${classes.icon}`}
              aria-hidden="true"
            >
              <Icon className="h-7 w-7 stroke-[2.5]" />
            </span>
            <div className="min-w-0">
              <span className="block truncate text-[12px] font-semibold leading-tight text-[#8f8798] sm:text-[13px]">
                {item.label}
              </span>
              <div className="relative mt-0.5 h-6 overflow-hidden">
                <AnimatePresence initial={false} mode="wait">
                  {isHidden ? (
                    <motion.span
                      key={`${item.label}-hidden`}
                      className="absolute left-0 top-1 block h-5 w-20 rounded-full bg-[#d7d1dd] sm:w-24"
                      aria-label={`${item.label} oculto`}
                      initial={{
                        opacity: 0,
                        y: shouldReduceMotion ? 0 : 3,
                        scale: shouldReduceMotion ? 1 : 0.98,
                      }}
                      animate={{ opacity: 1, y: 0, scale: 1 }}
                      exit={{
                        opacity: 0,
                        y: shouldReduceMotion ? 0 : -3,
                        scale: shouldReduceMotion ? 1 : 0.98,
                      }}
                      transition={valueTransition}
                    />
                  ) : (
                    <motion.strong
                      key={`${item.label}-visible`}
                      className={`absolute inset-x-0 top-0 block truncate text-[15px] font-semibold leading-tight sm:text-[16px] ${classes.value}`}
                      initial={{
                        opacity: 0,
                        y: shouldReduceMotion ? 0 : 3,
                        scale: shouldReduceMotion ? 1 : 0.98,
                      }}
                      animate={{ opacity: 1, y: 0, scale: 1 }}
                      exit={{
                        opacity: 0,
                        y: shouldReduceMotion ? 0 : -3,
                        scale: shouldReduceMotion ? 1 : 0.98,
                      }}
                      transition={valueTransition}
                    >
                      {formatSummaryValue(item.value)}
                    </motion.strong>
                  )}
                </AnimatePresence>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
