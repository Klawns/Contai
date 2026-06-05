import type { ReactNode } from 'react'
import { X } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'

type SelectionSheetProps = {
  title: string
  isOpen: boolean
  onClose: () => void
  children: ReactNode
}

export function SelectionSheet({ title, isOpen, onClose, children }: SelectionSheetProps) {
  const shouldReduceMotion = useReducedMotion()

  return (
    <AnimatePresence>
      {isOpen ? (
        <motion.div
          className="fixed inset-0 z-40 grid items-end bg-black/45 px-0 md:items-center md:px-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' }}
        >
          <button
            type="button"
            className="absolute inset-0 cursor-pointer"
            aria-label="Fechar selecao"
            onClick={onClose}
          />
          <motion.div
            role="dialog"
            aria-modal="true"
            aria-label={title}
            className="relative max-h-[82svh] overflow-hidden rounded-t-[22px] bg-white shadow-[0_-18px_48px_rgba(37,29,47,0.18)] md:mx-auto md:w-[min(100%,460px)] md:rounded-[18px]"
            initial={{ y: shouldReduceMotion ? 0 : 24, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: shouldReduceMotion ? 0 : 16, opacity: 0 }}
            transition={{ duration: shouldReduceMotion ? 0 : 0.18, ease: 'easeOut' }}
          >
            <div className="flex items-center justify-between border-b border-[#eee8f3] px-5 py-4">
              <h2 className="text-[16px] font-semibold text-[#2c2237]">{title}</h2>
              <button
                type="button"
                className="grid h-9 w-9 cursor-pointer place-items-center rounded-full text-[#81788c] transition-colors hover:bg-[#f7f2fb] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                aria-label="Fechar"
                onClick={onClose}
              >
                <X className="h-5 w-5" aria-hidden="true" />
              </button>
            </div>
            <div className="max-h-[calc(82svh-69px)] overflow-y-auto p-4">{children}</div>
          </motion.div>
        </motion.div>
      ) : null}
    </AnimatePresence>
  )
}
