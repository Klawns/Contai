import type { NavigationItem } from '../services/navigation'
import { motion, useReducedMotion } from 'motion/react'

type DockItemProps = {
  item: NavigationItem
  variant?: 'dock' | 'sidebar'
}

const baseButtonClasses =
  'min-w-0 cursor-pointer border-0 bg-transparent font-[inherit] text-[#b2aeb8] transition-[color,transform] duration-150 ease-in-out focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] hover:text-[#6a22e5]'

const dockClasses =
  'grid min-h-[54px] content-center justify-items-center gap-1 px-0.5 py-1 [&_svg]:h-[21px] [&_svg]:w-[21px] [&_span]:w-full [&_span]:overflow-hidden [&_span]:text-ellipsis [&_span]:whitespace-nowrap [&_span]:text-center [&_span]:text-[10px] [&_span]:font-semibold [&_span]:leading-[1.15] max-[360px]:[&_svg]:h-5 max-[360px]:[&_svg]:w-5 max-[360px]:[&_span]:text-[8px]'

const sidebarClasses =
  'grid min-h-11 grid-cols-[22px_minmax(0,1fr)] items-center justify-items-start rounded-lg px-2.5 py-[9px] text-[#86808f] [&_svg]:h-[19px] [&_svg]:w-[19px] [&_span]:w-full [&_span]:overflow-hidden [&_span]:text-ellipsis [&_span]:whitespace-nowrap [&_span]:text-left [&_span]:text-[13px] [&_span]:font-semibold [&_span]:leading-[1.15]'

export function DockItem({ item, variant = 'dock' }: DockItemProps) {
  const Icon = item.icon
  const shouldReduceMotion = useReducedMotion()
  const variantClasses = variant === 'sidebar' ? sidebarClasses : dockClasses
  const activeClasses =
    variant === 'sidebar' && item.active
      ? 'bg-[rgba(104,24,232,0.1)] text-[#6a22e5]'
      : item.active
        ? 'text-[#6a22e5]'
        : ''

  return (
    <motion.button
      type="button"
      className={`${baseButtonClasses} ${variantClasses} ${activeClasses}`}
      onClick={item.onSelect}
      whileHover={shouldReduceMotion ? undefined : { y: variant === 'dock' ? -1 : 0 }}
      whileTap={shouldReduceMotion ? undefined : { scale: 0.97 }}
    >
      <Icon />
      <span>{item.label}</span>
    </motion.button>
  )
}
