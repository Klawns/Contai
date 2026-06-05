import { UserRound } from 'lucide-react'
import { MonthSelector, type SelectedMonth } from '../../../../components/MonthSelector.tsx'

type DashboardHeaderProps = {
  selectedMonth: SelectedMonth
  userName: string
  onSelectMonth: (month: SelectedMonth) => void
}

export function DashboardHeader({
  selectedMonth,
  userName,
  onSelectMonth,
}: DashboardHeaderProps) {
  return (
    <header className="grid grid-cols-[44px_minmax(0,1fr)_44px] items-center pt-1">
      <button
        type="button"
        className="grid h-11 w-11 cursor-pointer place-items-center rounded-full border border-[#e4dfec] bg-white text-[#6b6178] shadow-[0_10px_24px_rgba(43,35,54,0.05)] transition-colors hover:text-[#6a22e5] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        aria-label={`Perfil de ${userName}`}
      >
        <UserRound className="h-5 w-5" aria-hidden="true" />
      </button>
      <div className="grid justify-center px-2">
        <MonthSelector
          selectedMonth={selectedMonth}
          onSelectMonth={onSelectMonth}
        />
      </div>
      <span aria-hidden="true" />
    </header>
  )
}
