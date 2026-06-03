import type { NavigationItem } from '../services/navigation'
import { Plus, X } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { DockItem } from './DockItem'

type BottomDockProps = {
  items: NavigationItem[]
  isQuickActionsOpen: boolean
  onToggleQuickActions: () => void
}

export function BottomDock({
  items,
  isQuickActionsOpen,
  onToggleQuickActions,
}: BottomDockProps) {
  const firstItems = items.slice(0, 2)
  const lastItems = items.slice(2)
  const QuickActionsIcon = isQuickActionsOpen ? X : Plus
  const shouldReduceMotion = useReducedMotion()

  return (
    <nav
      className="fixed inset-x-0 bottom-0 z-30 grid h-[calc(64px+env(safe-area-inset-bottom))] grid-cols-5 items-end border-t border-[#e8e3ef] bg-white px-[max(8px,env(safe-area-inset-left))] pb-[calc(6px+env(safe-area-inset-bottom))] pt-1.5 shadow-[0_-8px_24px_rgba(34,24,48,0.04)] md:hidden max-[360px]:px-1"
      aria-label="Navegacao principal"
    >
      {firstItems.map((item) => (
        <DockItem key={item.label} item={item} />
      ))}
      <motion.button
        type="button"
        className="place-self-start justify-self-center grid h-12 w-12 cursor-pointer place-items-center rounded-full border-0 bg-[#6818e8] font-[inherit] text-white shadow-[0_10px_22px_rgba(93,24,207,0.32)] transition-[box-shadow,transform] duration-150 ease-in-out hover:-translate-y-px hover:shadow-[0_12px_26px_rgba(93,24,207,0.42)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] -mt-5 [&_svg]:h-7 [&_svg]:w-7"
        aria-label={isQuickActionsOpen ? 'Fechar acoes rapidas' : 'Abrir acoes rapidas'}
        aria-expanded={isQuickActionsOpen}
        onClick={onToggleQuickActions}
        whileTap={shouldReduceMotion ? undefined : { scale: 0.94 }}
      >
        <AnimatePresence initial={false} mode="wait">
          <motion.span
            key={isQuickActionsOpen ? 'close' : 'open'}
            className="grid place-items-center"
            initial={{
              opacity: 0,
              rotate: shouldReduceMotion ? 0 : -45,
              scale: shouldReduceMotion ? 1 : 0.82,
            }}
            animate={{ opacity: 1, rotate: 0, scale: 1 }}
            exit={{
              opacity: 0,
              rotate: shouldReduceMotion ? 0 : 45,
              scale: shouldReduceMotion ? 1 : 0.82,
            }}
            transition={{ duration: shouldReduceMotion ? 0 : 0.14, ease: 'easeOut' }}
          >
            <QuickActionsIcon aria-hidden="true" />
          </motion.span>
        </AnimatePresence>
      </motion.button>
      {lastItems.map((item) => (
        <DockItem key={item.label} item={item} />
      ))}
    </nav>
  )
}
