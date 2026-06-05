import type { ReactNode } from 'react'
import { motion, useReducedMotion } from 'motion/react'

type DashboardLayoutProps = {
  animationKey?: string
  animateOnMount?: boolean
  children: ReactNode
  width?: 'contained' | 'full'
}

export function DashboardLayout({
  animationKey,
  animateOnMount = true,
  children,
  width = 'contained',
}: DashboardLayoutProps) {
  const shouldReduceMotion = useReducedMotion()
  const widthClasses =
    width === 'full'
      ? 'min-h-svh w-full max-w-none flex-1 gap-0 px-0 py-0'
      : 'mx-auto w-[min(100%,1120px)] gap-4 px-4 pb-[calc(92px+env(safe-area-inset-bottom))] pt-5 sm:px-5 md:gap-5 md:px-7 md:pb-9 md:pt-7 lg:px-10 lg:py-9'

  return (
    <motion.main
      key={animationKey}
      className={`grid text-left ${widthClasses}`}
      aria-label="Dashboard financeiro"
      initial={shouldReduceMotion || !animateOnMount ? false : { opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: shouldReduceMotion ? 0 : 0.24, ease: 'easeOut' }}
    >
      {children}
    </motion.main>
  )
}
