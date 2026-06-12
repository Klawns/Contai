import type { ReactNode } from 'react'
import { motion, useReducedMotion } from 'motion/react'

type TransactionsPageLayoutProps = {
  animationKey?: string
  variant?: 'default' | 'create'
  tone?: 'default' | 'income' | 'expense' | 'transfer'
  children: ReactNode
}

const toneClasses = {
  default: 'from-[#6f28e7] to-[#8c4df1]',
  income: 'from-[#159c57] to-[#28bd73]',
  expense: 'from-[#d93658] to-[#ef5b75]',
  transfer: 'from-[#2478d4] to-[#4aa8e8]',
}

export function TransactionsPageLayout({
  animationKey,
  variant = 'default',
  tone = 'default',
  children,
}: TransactionsPageLayoutProps) {
  const shouldReduceMotion = useReducedMotion()
  const isCreateVariant = variant === 'create'

  return (
    <motion.main
      key={animationKey}
      className={
        isCreateVariant
          ? `h-full min-h-0 w-full overflow-hidden bg-gradient-to-b ${toneClasses[tone]} from-0% via-[#f4f7fb] via-[215px] to-[#f4f7fb] text-left md:bg-[#f4f7fb] md:bg-none`
          : 'h-full min-h-0 w-full overflow-hidden bg-[#6818e8] text-left'
      }
      aria-label="Transacoes"
      initial={shouldReduceMotion ? false : { opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: shouldReduceMotion ? 0 : 0.22, ease: 'easeOut' }}
    >
      <div
        className={
          isCreateVariant
            ? 'grid h-full min-h-0 w-full min-w-0 overflow-hidden'
            : 'grid h-full min-h-0 w-full min-w-0 overflow-hidden'
        }
      >
        {children}
      </div>
    </motion.main>
  )
}
