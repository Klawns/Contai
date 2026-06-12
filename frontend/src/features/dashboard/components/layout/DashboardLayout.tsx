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
      ? 'h-full min-h-0 w-full max-w-none flex-1 gap-0 overflow-hidden px-0 py-0'
      : 'scrollbar-none mx-auto h-full min-h-0 w-[min(100%,1120px)] gap-4 overflow-y-auto overflow-x-hidden px-4 pb-[var(--app-mobile-content-bottom)] pt-5 sm:px-5 md:gap-5 md:px-7 md:pb-9 md:pt-7 lg:px-10 lg:py-9'

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
