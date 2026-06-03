import { useEffect, useRef, useState } from 'react'
import { UserRound, type LucideIcon } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'

export type ProfileAction = {
  label: string
  icon: LucideIcon
  tone?: 'default' | 'danger'
  disabled?: boolean
  onSelect: () => void
}

type ProfileActionsDropdownProps = {
  actions: ProfileAction[]
  ariaLabel: string
}

export function ProfileActionsDropdown({ actions, ariaLabel }: ProfileActionsDropdownProps) {
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const shouldReduceMotion = useReducedMotion()

  useEffect(() => {
    if (!isOpen) {
      return
    }

    const handlePointerDown = (event: PointerEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    window.addEventListener('pointerdown', handlePointerDown)

    return () => window.removeEventListener('pointerdown', handlePointerDown)
  }, [isOpen])

  return (
    <div ref={containerRef} className="absolute left-0 top-0 z-20">
      <button
        type="button"
        className="grid h-11 w-11 cursor-pointer place-items-center rounded-full border border-[#e4dfec] bg-white text-[#6b6178] shadow-[0_10px_24px_rgba(43,35,54,0.05)] transition-colors hover:text-[#6a22e5] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        aria-label={ariaLabel}
        aria-expanded={isOpen}
        aria-haspopup="menu"
        onClick={() => setIsOpen((current) => !current)}
      >
        <UserRound className="h-5 w-5" aria-hidden="true" />
      </button>

      <AnimatePresence>
        {isOpen && (
          <motion.div
            className="absolute left-0 top-[52px] min-w-[156px] rounded-[12px] border border-[#ece8f2] bg-white p-1.5 text-left shadow-[0_16px_38px_rgba(48,39,61,0.14)]"
            role="menu"
            aria-label="Acoes do perfil"
            initial={{
              opacity: 0,
              y: shouldReduceMotion ? 0 : -4,
              scale: shouldReduceMotion ? 1 : 0.96,
            }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{
              opacity: 0,
              y: shouldReduceMotion ? 0 : -4,
              scale: shouldReduceMotion ? 1 : 0.96,
            }}
            style={{ transformOrigin: 'top left' }}
            transition={{ duration: shouldReduceMotion ? 0 : 0.14, ease: 'easeOut' }}
          >
            {actions.map((action) => {
              const Icon = action.icon
              const isDanger = action.tone === 'danger'

              return (
                <button
                  key={action.label}
                  type="button"
                  className={`flex h-10 w-full cursor-pointer items-center gap-2 rounded-[9px] px-2.5 text-[13px] font-semibold transition-colors disabled:cursor-not-allowed ${
                    isDanger
                      ? 'text-[#c93434] hover:bg-[#fff4f4] disabled:text-[#d58d8d]'
                      : 'text-[#4d4655] hover:bg-[#f6f2fb] disabled:text-[#9f97aa]'
                  }`}
                  role="menuitem"
                  disabled={action.disabled}
                  onClick={() => {
                    action.onSelect()
                    setIsOpen(false)
                  }}
                >
                  <Icon className="h-4 w-4" aria-hidden="true" />
                  <span>{action.label}</span>
                </button>
              )
            })}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
