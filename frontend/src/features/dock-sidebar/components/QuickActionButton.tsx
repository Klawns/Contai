import type { CSSProperties } from 'react'
import { motion, useReducedMotion } from 'motion/react'
import type { QuickAction } from '../services/navigation'

type QuickActionButtonProps = {
  action: QuickAction
  index?: number
  className?: string
}

export function QuickActionButton({
  action,
  index = 0,
  className = '',
}: QuickActionButtonProps) {
  const Icon = action.icon
  const style = { '--action-color': action.color } as CSSProperties
  const shouldReduceMotion = useReducedMotion()

  return (
    <motion.button
      type="button"
      className={`grid min-w-0 cursor-pointer justify-items-center gap-2 border-0 bg-transparent p-0 font-[inherit] text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${className}`}
      style={style}
      onClick={action.onSelect}
      initial={{
        opacity: 0,
        y: shouldReduceMotion ? 0 : 10,
        scale: shouldReduceMotion ? 1 : 0.94,
      }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{
        opacity: 0,
        y: shouldReduceMotion ? 0 : 6,
        scale: shouldReduceMotion ? 1 : 0.96,
      }}
      transition={{
        duration: shouldReduceMotion ? 0 : 0.16,
        delay: shouldReduceMotion ? 0 : index * 0.035,
        ease: 'easeOut',
      }}
      whileTap={shouldReduceMotion ? undefined : { scale: 0.96 }}
    >
      <span className="grid h-14 w-14 place-items-center rounded-full bg-white text-[var(--action-color)] shadow-[0_12px_26px_rgba(0,0,0,0.22)] [&_svg]:h-[26px] [&_svg]:w-[26px]">
        <Icon />
      </span>
      <span className="block w-[104px] text-center text-xs font-semibold leading-[1.15] text-white [overflow-wrap:anywhere] max-[360px]:w-[92px] max-[360px]:text-[11px]">
        {action.label}
      </span>
    </motion.button>
  )
}
