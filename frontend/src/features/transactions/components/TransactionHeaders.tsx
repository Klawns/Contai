import { ArrowLeft, Plus } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import type { AuthPath } from '../../auth/services/navigation.ts'
import { MonthSelector, type SelectedMonth } from '../../../components/MonthSelector.tsx'

type TransactionMobileHeaderProps = {
  title: string
  selectedMonth?: SelectedMonth
  onSelectMonth?: (month: SelectedMonth) => void
}

export function TransactionMobileHeader({
  title,
  selectedMonth,
  onSelectMonth,
}: TransactionMobileHeaderProps) {
  const navigate = useNavigate()

  return (
    <header className="grid grid-cols-[44px_minmax(0,1fr)_44px] items-center text-white md:hidden">
      <button
        type="button"
        className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
        aria-label="Voltar"
        onClick={() => navigate('/')}
      >
        <ArrowLeft className="h-5 w-5" aria-hidden="true" />
      </button>
      <div className="min-w-0 px-2 text-center">
        <h1 className="truncate text-[17px] font-semibold leading-tight">{title}</h1>
        {selectedMonth && onSelectMonth ? (
          <div className="mt-1 flex justify-center [&_button]:text-white [&_span]:text-white">
            <MonthSelector selectedMonth={selectedMonth} onSelectMonth={onSelectMonth} />
          </div>
        ) : null}
      </div>
      <button
        type="button"
        className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white text-[#6a22e5] shadow-[0_10px_24px_rgba(43,35,54,0.12)] transition-colors hover:bg-[#f7f2ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
        aria-label="Nova despesa"
        onClick={() => navigate('/transactions/expense/new')}
      >
        <Plus className="h-5 w-5" aria-hidden="true" />
      </button>
    </header>
  )
}

type TransactionDesktopHeaderProps = {
  title: string
  subtitle: string
  selectedMonth?: SelectedMonth
  onSelectMonth?: (month: SelectedMonth) => void
  actionPath?: AuthPath
}

export function TransactionDesktopHeader({
  title,
  subtitle,
  selectedMonth,
  onSelectMonth,
  actionPath = '/transactions/expense/new',
}: TransactionDesktopHeaderProps) {
  const navigate = useNavigate()

  return (
    <header className="hidden items-center justify-between gap-4 md:flex">
      <div className="min-w-0">
        <p className="text-[13px] font-semibold uppercase tracking-[0.08em] text-[#7d718b]">
          Contai
        </p>
        <h1 className="mt-1 text-[26px] font-semibold leading-tight text-[#281d35]">
          {title}
        </h1>
        <p className="mt-1 text-[14px] font-medium text-[#81788c]">{subtitle}</p>
      </div>
      <div className="flex items-center gap-3">
        {selectedMonth && onSelectMonth ? (
          <div className="rounded-full border border-[#e5deef] bg-white px-4 shadow-[0_10px_24px_rgba(43,35,54,0.05)]">
            <MonthSelector selectedMonth={selectedMonth} onSelectMonth={onSelectMonth} />
          </div>
        ) : null}
        <button
          type="button"
          className="inline-flex h-11 cursor-pointer items-center gap-2 rounded-lg bg-[#6818e8] px-4 text-[14px] font-semibold text-white shadow-[0_10px_22px_rgba(93,24,207,0.24)] transition-colors hover:bg-[#5812cf] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
          onClick={() => navigate(actionPath)}
        >
          <Plus className="h-4 w-4" aria-hidden="true" />
          Nova transacao
        </button>
      </div>
    </header>
  )
}
