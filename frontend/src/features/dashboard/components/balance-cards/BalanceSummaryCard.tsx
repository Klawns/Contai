import { Eye, EyeOff } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import {
  ProfileActionsDropdown,
  type ProfileAction,
} from '../../../auth/components/ProfileActionsDropdown.tsx'
import { MonthSelector, type SelectedMonth } from '../../../../components/MonthSelector.tsx'
import { IncomeExpenseSummary } from './IncomeExpenseSummary.tsx'

type BalanceSummaryCardProps = {
  selectedMonth: SelectedMonth
  userName: string
  totalBalance: number
  monthlyIncome: number
  monthlyExpense: number
  isBalanceHidden?: boolean
  profileActions: ProfileAction[]
  onSelectMonth: (month: SelectedMonth) => void
  onToggleVisibility: () => void
}

const mainBalanceFormatter = new Intl.NumberFormat('pt-BR', {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

function formatMainBalance(valueInCents: number) {
  const absolute = Math.abs(valueInCents) / 100
  const formatted = mainBalanceFormatter.format(absolute)

  return valueInCents < 0 ? `-${formatted}` : formatted
}

export function BalanceSummaryCard({
  selectedMonth,
  userName,
  totalBalance,
  monthlyIncome,
  monthlyExpense,
  isBalanceHidden = false,
  profileActions,
  onSelectMonth,
  onToggleVisibility,
}: BalanceSummaryCardProps) {
  const shouldReduceMotion = useReducedMotion()
  const VisibilityIcon = isBalanceHidden ? EyeOff : Eye
  const valueTransition = { duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' as const }
  const iconTransition = { duration: shouldReduceMotion ? 0 : 0.14, ease: 'easeOut' as const }

  return (
    <motion.section
      className="-mx-4 -mt-5 rounded-b-[24px] border border-x-0 border-t-0 border-[#ece8f2] bg-white px-5 pb-6 pt-4 shadow-[0_18px_44px_rgba(48,39,61,0.08)] sm:mx-0 sm:mt-0 sm:rounded-[22px] sm:border sm:px-8 sm:py-7"
      aria-label="Saldo atual em contas"
      whileHover={shouldReduceMotion ? undefined : { y: -2 }}
      transition={{ duration: 0.18, ease: 'easeOut' }}
    >
      <div className="mb-7 grid justify-items-center gap-4 text-center sm:gap-5">
        <div className="relative grid w-full place-items-center md:hidden">
          <ProfileActionsDropdown
            actions={profileActions}
            ariaLabel={`Perfil de ${userName}`}
          />
          <div className="grid max-w-[calc(100%-104px)] justify-center">
            <MonthSelector
              selectedMonth={selectedMonth}
              onSelectMonth={onSelectMonth}
            />
          </div>
        </div>

        <div className="hidden md:block">
          <MonthSelector
            selectedMonth={selectedMonth}
            onSelectMonth={onSelectMonth}
          />
        </div>

        <div className="grid w-full justify-items-center gap-2">
          <span className="text-[12px] font-medium leading-none tracking-0 text-[#8f8798] sm:text-[13px]">
            Saldo atual em contas
          </span>
          <div className="relative grid h-10 w-full place-items-center overflow-hidden sm:h-12">
            <AnimatePresence initial={false} mode="wait">
              {isBalanceHidden ? (
                <motion.span
                  key="hidden-balance"
                  className="block h-10 w-[230px] max-w-[78vw] rounded-full bg-[#d7d1dd] sm:h-12 sm:w-[292px]"
                  aria-label="Saldo oculto"
                  initial={{
                    opacity: 0,
                    y: shouldReduceMotion ? 0 : 4,
                    scale: shouldReduceMotion ? 1 : 0.98,
                  }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  exit={{
                    opacity: 0,
                    y: shouldReduceMotion ? 0 : -4,
                    scale: shouldReduceMotion ? 1 : 0.98,
                  }}
                  transition={valueTransition}
                />
              ) : (
                <motion.strong
                  key="visible-balance"
                  className="block max-w-full overflow-hidden text-ellipsis whitespace-nowrap text-[36px] font-semibold leading-none tracking-0 text-[#211827] sm:text-[46px]"
                  initial={{
                    opacity: 0,
                    y: shouldReduceMotion ? 0 : 4,
                    scale: shouldReduceMotion ? 1 : 0.98,
                  }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  exit={{
                    opacity: 0,
                    y: shouldReduceMotion ? 0 : -4,
                    scale: shouldReduceMotion ? 1 : 0.98,
                  }}
                  transition={valueTransition}
                >
                  {formatMainBalance(totalBalance)}
                </motion.strong>
              )}
            </AnimatePresence>
          </div>
        </div>

        <motion.button
          type="button"
          className="grid h-8 w-8 cursor-pointer place-items-center text-[#4d4655] transition-colors hover:text-[#6a22e5] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
          aria-label={isBalanceHidden ? 'Mostrar valores do saldo' : 'Ocultar valores do saldo'}
          aria-pressed={isBalanceHidden}
          onClick={onToggleVisibility}
          whileTap={shouldReduceMotion ? undefined : { scale: 0.92 }}
        >
          <AnimatePresence initial={false} mode="wait">
            <motion.span
              key={isBalanceHidden ? 'hidden-icon' : 'visible-icon'}
              className="grid place-items-center"
              initial={{
                opacity: 0,
                rotate: shouldReduceMotion ? 0 : -20,
                scale: shouldReduceMotion ? 1 : 0.86,
              }}
              animate={{ opacity: 1, rotate: 0, scale: 1 }}
              exit={{
                opacity: 0,
                rotate: shouldReduceMotion ? 0 : 20,
                scale: shouldReduceMotion ? 1 : 0.86,
              }}
              transition={iconTransition}
            >
              <VisibilityIcon className="h-5 w-5" aria-hidden="true" />
            </motion.span>
          </AnimatePresence>
        </motion.button>
      </div>

      <IncomeExpenseSummary
        income={monthlyIncome}
        expense={monthlyExpense}
        isHidden={isBalanceHidden}
      />
    </motion.section>
  )
}
