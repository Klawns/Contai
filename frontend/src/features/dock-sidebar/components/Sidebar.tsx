import { useState } from 'react'
import { ChevronDown, UserRound } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { LogoutActionButton } from '../../auth/components/LogoutActionButton'
import type { NavigationItem } from '../services/navigation'
import { DockItem } from './DockItem'

type SidebarProps = {
  items: NavigationItem[]
  isLoggingOut: boolean
  onLogout: () => void
}

const parentButtonClasses =
  'grid min-h-11 w-full cursor-pointer grid-cols-[22px_minmax(0,1fr)_16px] items-center justify-items-start rounded-lg border-0 bg-transparent px-2.5 py-[9px] font-[inherit] text-[#86808f] transition-[background-color,color,transform] duration-150 ease-in-out hover:text-[#6a22e5] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] [&_svg]:h-[19px] [&_svg]:w-[19px] [&_span]:w-full [&_span]:overflow-hidden [&_span]:text-ellipsis [&_span]:whitespace-nowrap [&_span]:text-left [&_span]:text-[13px] [&_span]:font-semibold [&_span]:leading-[1.15]'

const parentActiveClasses = 'bg-[rgba(104,24,232,0.1)] text-[#6a22e5]'

function getInitialOpenItems(items: NavigationItem[]) {
  return new Set(items.filter((item) => item.children && item.active).map((item) => item.label))
}

export function Sidebar({ items, isLoggingOut, onLogout }: SidebarProps) {
  const [openItems, setOpenItems] = useState(() => getInitialOpenItems(items))
  const [closedItems, setClosedItems] = useState<Set<string>>(() => new Set())
  const shouldReduceMotion = useReducedMotion()

  function toggleItem(label: string, isOpen: boolean) {
    setOpenItems((currentOpenItems) => {
      const nextOpenItems = new Set(currentOpenItems)

      if (isOpen) {
        nextOpenItems.delete(label)
      } else {
        nextOpenItems.add(label)
      }

      return nextOpenItems
    })

    setClosedItems((currentClosedItems) => {
      const nextClosedItems = new Set(currentClosedItems)

      if (isOpen) {
        nextClosedItems.add(label)
      } else {
        nextClosedItems.delete(label)
      }

      return nextClosedItems
    })
  }

  return (
    <aside
      className="sticky top-0 hidden h-svh flex-col gap-6 border-r border-[#e6e0ee] bg-white px-4 py-6 md:flex"
      aria-label="Navegacao principal"
    >
      <div className="flex items-center gap-2.5 text-lg text-[#241932]">
        <span className="grid h-[34px] w-[34px] place-items-center rounded-full border border-[#e4dfec] bg-[#fbfafe] text-[#6b6178]">
          <UserRound className="h-5 w-5" aria-hidden="true" />
        </span>
        <strong className="font-medium">Contai</strong>
      </div>
      <nav className="grid gap-1.5">
        {items.map((item) => {
          if (!item.children) {
            return <DockItem key={item.label} item={item} variant="sidebar" />
          }

          const Icon = item.icon
          const isOpen =
            openItems.has(item.label) || Boolean(item.active && !closedItems.has(item.label))

          return (
            <div key={item.label} className="grid gap-1">
              <motion.button
                type="button"
                className={`${parentButtonClasses} ${item.active ? parentActiveClasses : ''}`}
                aria-expanded={isOpen}
                onClick={() => toggleItem(item.label, isOpen)}
                whileTap={shouldReduceMotion ? undefined : { scale: 0.97 }}
              >
                <Icon aria-hidden="true" />
                <span>{item.label}</span>
                <ChevronDown
                  className={`justify-self-end transition-transform duration-150 ${
                    isOpen ? 'rotate-180' : ''
                  }`}
                  aria-hidden="true"
                />
              </motion.button>
              <AnimatePresence initial={false}>
                {isOpen ? (
                  <motion.div
                    className="grid overflow-hidden"
                    initial={{
                      height: 0,
                      opacity: 0,
                      x: shouldReduceMotion ? 0 : -4,
                    }}
                    animate={{ height: 'auto', opacity: 1, x: 0 }}
                    exit={{
                      height: 0,
                      opacity: 0,
                      x: shouldReduceMotion ? 0 : -4,
                    }}
                    transition={{ duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' }}
                  >
                    <div className="grid gap-1 pb-1 pl-[34px] pt-0.5">
                      {item.children.map((child) => (
                        <button
                          key={child.path}
                          type="button"
                          className={`min-h-9 cursor-pointer rounded-md border-0 bg-transparent px-2.5 py-2 text-left font-[inherit] text-[12px] font-semibold leading-[1.15] text-[#8f879a] transition-[background-color,color] duration-150 ease-in-out hover:bg-[#f7f3fd] hover:text-[#6a22e5] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
                            child.active ? 'bg-[#f2ebff] text-[#6a22e5]' : ''
                          }`}
                          onClick={child.onSelect}
                        >
                          {child.label}
                        </button>
                      ))}
                    </div>
                  </motion.div>
                ) : null}
              </AnimatePresence>
            </div>
          )
        })}
      </nav>
      <div className="mt-auto">
        <LogoutActionButton
          variant="sidebar"
          isLoggingOut={isLoggingOut}
          onLogout={onLogout}
        />
      </div>
    </aside>
  )
}
