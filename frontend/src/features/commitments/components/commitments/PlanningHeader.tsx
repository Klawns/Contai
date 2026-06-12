import { ArrowLeft, ChevronLeft, ChevronRight, Plus } from 'lucide-react'
import { MonthSelector, type SelectedMonth } from '../../../../components/MonthSelector.tsx'

type PlanningHeaderProps = {
  selectedMonth: SelectedMonth
  onBack: () => void
  onCreate: () => void
  onChangeMonth: (offset: -1 | 1) => void
  onSelectMonth: (month: SelectedMonth) => void
}

export function PlanningHeader({
  selectedMonth,
  onBack,
  onCreate,
  onChangeMonth,
  onSelectMonth,
}: PlanningHeaderProps) {
  return (
    <header className="flex-none bg-[#6818e8] px-5 pb-5 pt-[calc(18px+env(safe-area-inset-top))] text-white md:px-7 md:pt-6">
      <div className="mx-auto grid w-full grid-cols-[44px_minmax(0,1fr)_44px] items-center">
        <button
          type="button"
          className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
          aria-label="Voltar"
          onClick={onBack}
        >
          <ArrowLeft className="h-5 w-5" aria-hidden="true" />
        </button>
        <div className="min-w-0 px-2 text-center">
          <h1 className="truncate text-[17px] font-semibold leading-tight md:text-[24px]">
            Planejamento
          </h1>
          <div className="mt-1 flex min-w-0 items-center justify-center gap-1">
            <button
              type="button"
              className="grid h-8 w-8 flex-none cursor-pointer place-items-center rounded-full text-white/88 transition-colors hover:bg-white/12 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              aria-label="Mes anterior"
              onClick={() => onChangeMonth(-1)}
            >
              <ChevronLeft className="h-4.5 w-4.5" aria-hidden="true" />
            </button>
            <MonthSelector selectedMonth={selectedMonth} onSelectMonth={onSelectMonth} tone="inverse" />
            <button
              type="button"
              className="grid h-8 w-8 flex-none cursor-pointer place-items-center rounded-full text-white/88 transition-colors hover:bg-white/12 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              aria-label="Proximo mes"
              onClick={() => onChangeMonth(1)}
            >
              <ChevronRight className="h-4.5 w-4.5" aria-hidden="true" />
            </button>
          </div>
        </div>
        <button
          type="button"
          className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
          aria-label="Novo compromisso"
          onClick={onCreate}
        >
          <Plus className="h-5 w-5" aria-hidden="true" />
        </button>
      </div>
    </header>
  )
}
