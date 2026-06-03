import type { ReactNode } from 'react'
import { motion, useReducedMotion } from 'motion/react'

type DashboardLayoutProps = {
  animationKey?: string
  animateOnMount?: boolean
  children: ReactNode
}

export function DashboardLayout({
  animationKey,
  animateOnMount = true,
  children,
}: DashboardLayoutProps) {
  const shouldReduceMotion = useReducedMotion()

  return (
    <motion.main
      key={animationKey}
      className="mx-auto grid w-[min(100%,1120px)] gap-4 px-4 pb-[calc(92px+env(safe-area-inset-bottom))] pt-5 text-left sm:px-5 md:gap-5 md:px-7 md:pb-9 md:pt-7 lg:px-10 lg:py-9"
      aria-label="Dashboard financeiro"
      initial={shouldReduceMotion || !animateOnMount ? false : { opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: shouldReduceMotion ? 0 : 0.24, ease: 'easeOut' }}
    >
      {children}
    </motion.main>
  )
}
