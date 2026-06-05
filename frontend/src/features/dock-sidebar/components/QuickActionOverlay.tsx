import type { QuickAction } from '../services/navigation'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { QuickActionButton } from './QuickActionButton'

type QuickActionOverlayProps = {
  actions: QuickAction[]
  isOpen: boolean
  onClose: () => void
}

export function QuickActionOverlay({
  actions,
  isOpen,
  onClose,
}: QuickActionOverlayProps) {
  const shouldReduceMotion = useReducedMotion()

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          className="pointer-events-auto fixed inset-0 z-20 grid items-end md:hidden"
          aria-hidden={!isOpen}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' }}
        >
          <motion.button
            type="button"
            className="absolute inset-0 cursor-pointer border-0 bg-black/70 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
            aria-label="Fechar acoes rapidas"
            onClick={onClose}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' }}
          />
          <motion.div
            className="relative mx-auto grid w-full max-w-[280px] grid-cols-2 justify-items-center gap-x-8 gap-y-5 px-5 pb-[calc(98px+env(safe-area-inset-bottom))] max-[360px]:max-w-[244px] max-[360px]:gap-x-5 max-[360px]:gap-y-4 max-[360px]:px-4"
            aria-label="Acoes rapidas"
            initial={{
              opacity: 0,
              y: shouldReduceMotion ? 0 : 14,
              scale: shouldReduceMotion ? 1 : 0.96,
            }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{
              opacity: 0,
              y: shouldReduceMotion ? 0 : 10,
              scale: shouldReduceMotion ? 1 : 0.98,
            }}
            transition={{ duration: shouldReduceMotion ? 0 : 0.18, ease: 'easeOut' }}
          >
            {actions.map((action, index) => (
              <QuickActionButton
                key={action.label}
                action={{
                  ...action,
                  onSelect: () => {
                    action.onSelect?.()
                    onClose()
                  },
                }}
                index={index}
                className={index === 0 ? 'col-span-2' : ''}
              />
            ))}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
