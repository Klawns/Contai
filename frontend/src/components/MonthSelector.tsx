import { useEffect, useRef, useState } from 'react'
import { ChevronDown, ChevronLeft, ChevronRight } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'

export type SelectedMonth = {
  year: number
  monthIndex: number
}

type MonthSelectorTone = 'default' | 'inverse'

type MonthSelectorProps = {
  selectedMonth: SelectedMonth
  onSelectMonth: (month: SelectedMonth) => void
  tone?: MonthSelectorTone
}

const monthFormatter = new Intl.DateTimeFormat('pt-BR', {
  month: 'long',
})

const monthButtonFormatter = new Intl.DateTimeFormat('pt-BR', {
  month: 'short',
})

const months = Array.from({ length: 12 }, (_, monthIndex) => ({
  monthIndex,
  label: monthButtonFormatter
    .format(new Date(2026, monthIndex, 1))
    .replace('.', ''),
}))

function getCapitalizedMonthName(monthIndex: number) {
  const label = monthFormatter.format(new Date(2026, monthIndex, 1))

  return label.charAt(0).toUpperCase() + label.slice(1)
}

const toneStyles: Record<
  MonthSelectorTone,
  {
    button: string
    label: string
  }
> = {
  default: {
    button:
      'text-[#8d8495] hover:text-[#5f566b]',
    label: 'text-[#5f566b]',
  },
  inverse: {
    button:
      'text-white/88 hover:text-white',
    label: 'text-white',
  },
}

export function MonthSelector({
  selectedMonth,
  onSelectMonth,
  tone = 'default',
}: MonthSelectorProps) {
  const shouldReduceMotion = useReducedMotion()
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const selectedMonthLabel = getCapitalizedMonthName(selectedMonth.monthIndex)
  const selectedTone = toneStyles[tone]

  useEffect(() => {
    if (!isOpen) {
      return undefined
    }

    function handlePointerDown(event: PointerEvent) {
      if (!containerRef.current?.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('pointerdown', handlePointerDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen])

  return (
    <div ref={containerRef} className="relative inline-flex min-w-0 items-center justify-center">
      <motion.button
        type="button"
        className={`inline-flex h-9 min-w-0 cursor-pointer items-center justify-center gap-1 bg-transparent px-1 text-[15px] font-medium leading-none transition-colors focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7b2cff] ${selectedTone.button}`}
        aria-label={`Mes selecionado: ${selectedMonthLabel}`}
        aria-expanded={isOpen}
        aria-haspopup="dialog"
        onClick={() => setIsOpen((current) => !current)}
        whileTap={shouldReduceMotion ? undefined : { scale: 0.97 }}
      >
        <span className={`truncate ${selectedTone.label}`}>{selectedMonthLabel}</span>
        <ChevronDown
          className={`relative top-px h-4 w-4 flex-none transition-transform ${isOpen ? 'rotate-180' : ''}`}
          aria-hidden="true"
        />
      </motion.button>

      {isOpen ? (
        <div
          role="dialog"
          aria-label="Selecionar mes"
          className="absolute left-1/2 top-11 z-30 w-[276px] -translate-x-1/2 rounded-2xl border border-[#e9e2f0] bg-white p-3 text-left shadow-[0_18px_48px_rgba(37,29,47,0.14)]"
        >
          <div className="mb-3 flex items-center justify-between">
            <button
              type="button"
              className="grid h-8 w-8 cursor-pointer place-items-center rounded-full text-[#81798b] transition-colors hover:bg-[#f6f2fb] hover:text-[#5f566b] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
              aria-label="Ano anterior"
              onClick={() => onSelectMonth({ ...selectedMonth, year: selectedMonth.year - 1 })}
            >
              <ChevronLeft className="h-4 w-4" aria-hidden="true" />
            </button>
            <span className="text-[14px] font-semibold leading-none text-[#3f354a]">
              {selectedMonth.year}
            </span>
            <button
              type="button"
              className="grid h-8 w-8 cursor-pointer place-items-center rounded-full text-[#81798b] transition-colors hover:bg-[#f6f2fb] hover:text-[#5f566b] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
              aria-label="Proximo ano"
              onClick={() => onSelectMonth({ ...selectedMonth, year: selectedMonth.year + 1 })}
            >
              <ChevronRight className="h-4 w-4" aria-hidden="true" />
            </button>
          </div>

          <div className="grid grid-cols-3 gap-1.5">
            {months.map((month) => {
              const isSelected = month.monthIndex === selectedMonth.monthIndex

              return (
                <button
                  key={month.monthIndex}
                  type="button"
                  className={`h-9 cursor-pointer rounded-lg px-2 text-[13px] font-medium capitalize transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
                    isSelected
                      ? 'bg-[#7b2cff] text-white shadow-[0_8px_18px_rgba(123,44,255,0.25)]'
                      : 'text-[#6f657a] hover:bg-[#f6f2fb] hover:text-[#3f354a]'
                  }`}
                  aria-pressed={isSelected}
                  onClick={() => {
                    onSelectMonth({ ...selectedMonth, monthIndex: month.monthIndex })
                    setIsOpen(false)
                  }}
                >
                  {month.label}
                </button>
              )
            })}
          </div>
        </div>
      ) : null}
    </div>
  )
}
