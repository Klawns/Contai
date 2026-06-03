import { useCallback, useEffect, useMemo, useState } from 'react'
import { LogOut } from 'lucide-react'
import type { AuthenticatedUser } from '../auth/types/auth.ts'
import type { ProfileAction } from '../auth/components/ProfileActionsDropdown.tsx'
import {
  AccountListCard,
  BalanceSummaryCard,
  DashboardChartsCarousel,
  DashboardLayout,
  DashboardSection,
  MonthlyBalanceCard,
} from './components'
import { summaryDisplayState } from './data/dashboardMocks.ts'
import { useMonthlyDashboard, useMonthlySeries } from './hooks/useMonthlyDashboard.ts'
import type { DashboardPeriod, MonthlyFinancialSeriesPoint } from './types/dashboard.ts'
import type { SelectedMonth } from './components/controls/MonthSelector.tsx'

type DashboardProps = {
  user: AuthenticatedUser
  isLoggingOut: boolean
  onLogout: () => void
}

const monthFormatter = new Intl.DateTimeFormat('pt-BR', {
  month: 'long',
})

const monthQueryPattern = /^(\d{4})-(0[1-9]|1[0-2])$/

function formatLocalRFC3339(date: Date) {
  const offsetInMinutes = -date.getTimezoneOffset()
  const offsetSign = offsetInMinutes >= 0 ? '+' : '-'
  const absoluteOffset = Math.abs(offsetInMinutes)
  const offsetHours = String(Math.floor(absoluteOffset / 60)).padStart(2, '0')
  const offsetMinutes = String(absoluteOffset % 60).padStart(2, '0')
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${offsetSign}${offsetHours}:${offsetMinutes}`
}

function getCurrentSelectedMonth(): SelectedMonth {
  const now = new Date()

  return {
    year: now.getFullYear(),
    monthIndex: now.getMonth(),
  }
}

function parseMonthQuery(value: string | null): SelectedMonth | null {
  const match = value?.match(monthQueryPattern)

  if (!match) {
    return null
  }

  return {
    year: Number(match[1]),
    monthIndex: Number(match[2]) - 1,
  }
}

function getInitialSelectedMonth(): SelectedMonth {
  if (typeof window === 'undefined') {
    return getCurrentSelectedMonth()
  }

  const params = new URLSearchParams(window.location.search)

  return parseMonthQuery(params.get('month')) ?? getCurrentSelectedMonth()
}

function formatMonthQuery(month: SelectedMonth) {
  return `${month.year}-${String(month.monthIndex + 1).padStart(2, '0')}`
}

function getMonthPeriod(month: SelectedMonth): DashboardPeriod {
  const startAt = new Date(month.year, month.monthIndex, 1, 0, 0, 0)
  const endAt = new Date(month.year, month.monthIndex + 1, 0, 23, 59, 59)

  return {
    startAt: formatLocalRFC3339(startAt),
    endAt: formatLocalRFC3339(endAt),
  }
}

function getMonthlySeriesPeriod(month: SelectedMonth): DashboardPeriod {
  const startAt = new Date(month.year, month.monthIndex - 5, 1, 0, 0, 0)
  const endAt = new Date(month.year, month.monthIndex + 1, 0, 23, 59, 59)

  return {
    startAt: formatLocalRFC3339(startAt),
    endAt: formatLocalRFC3339(endAt),
  }
}

function getMonthLabel(startAt: string) {
  const label = monthFormatter.format(new Date(startAt))

  return label.charAt(0).toUpperCase() + label.slice(1)
}

function toChartSeries(points: { monthStartAt: string; income: number; expense: number; balance: number }[]): MonthlyFinancialSeriesPoint[] {
  return points.map((point) => ({
    monthLabel: getMonthLabel(point.monthStartAt),
    income: point.income,
    expense: point.expense,
    balance: point.balance,
  }))
}

export function Dashboard({ user, isLoggingOut, onLogout }: DashboardProps) {
  const [selectedMonth, setSelectedMonth] = useState(getInitialSelectedMonth)
  const period = useMemo(() => getMonthPeriod(selectedMonth), [selectedMonth])
  const seriesPeriod = useMemo(() => getMonthlySeriesPeriod(selectedMonth), [selectedMonth])
  const { data: monthlyDashboard, isLoading, isError } = useMonthlyDashboard(period)
  const { data: monthlySeries } = useMonthlySeries(seriesPeriod)
  const chartSeries = useMemo(
    () => toChartSeries(monthlySeries?.points ?? []),
    [monthlySeries?.points],
  )
  const [isMainBalanceHidden, setIsMainBalanceHidden] = useState(
    summaryDisplayState.isBalanceHidden,
  )
  const updateSelectedMonth = useCallback((nextMonth: SelectedMonth) => {
    setSelectedMonth(nextMonth)

    if (typeof window === 'undefined') {
      return
    }

    const url = new URL(window.location.href)
    url.searchParams.set('month', formatMonthQuery(nextMonth))
    window.history.pushState(null, '', `${url.pathname}${url.search}${url.hash}`)
  }, [])
  const profileActions: ProfileAction[] = [
    {
      label: isLoggingOut ? 'Saindo...' : 'Sair',
      icon: LogOut,
      tone: 'danger',
      disabled: isLoggingOut,
      onSelect: onLogout,
    },
  ]

  useEffect(() => {
    function handlePopState() {
      const params = new URLSearchParams(window.location.search)
      setSelectedMonth(parseMonthQuery(params.get('month')) ?? getCurrentSelectedMonth())
    }

    window.addEventListener('popstate', handlePopState)

    return () => {
      window.removeEventListener('popstate', handlePopState)
    }
  }, [])

  if (isLoading) {
    return (
      <DashboardLayout animateOnMount={false}>
        <div className="rounded-[18px] border border-[#ece8f2] bg-white px-5 py-6 text-[14px] font-medium text-[#81798b] shadow-[0_16px_38px_rgba(48,39,61,0.07)]">
          Carregando dashboard...
        </div>
      </DashboardLayout>
    )
  }

  if (isError || !monthlyDashboard) {
    return (
      <DashboardLayout>
        <div className="rounded-[18px] border border-[#f0caca] bg-white px-5 py-6 text-[14px] font-medium text-[#b93838] shadow-[0_16px_38px_rgba(48,39,61,0.07)]">
          Nao foi possivel carregar o dashboard.
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout animationKey={formatMonthQuery(selectedMonth)}>
      <BalanceSummaryCard
        selectedMonth={selectedMonth}
        userName={user.name ?? user.email}
        totalBalance={monthlyDashboard.totalBalance}
        monthlyIncome={monthlyDashboard.monthlyIncome}
        monthlyExpense={monthlyDashboard.monthlyExpense}
        isBalanceHidden={isMainBalanceHidden}
        profileActions={profileActions}
        onSelectMonth={updateSelectedMonth}
        onToggleVisibility={() => setIsMainBalanceHidden((isHidden) => !isHidden)}
      />

      <DashboardSection title="Contas">
        <AccountListCard
          accounts={monthlyDashboard.accountBalances}
          isBalanceHidden={isMainBalanceHidden}
        />
      </DashboardSection>

      <DashboardChartsCarousel
        expensesByCategory={monthlyDashboard.expensesByCategory}
        monthlySeries={chartSeries}
      >
        <MonthlyBalanceCard
          monthlyIncome={monthlyDashboard.monthlyIncome}
          monthlyExpense={monthlyDashboard.monthlyExpense}
          monthlyNetBalance={monthlyDashboard.monthlyNetBalance}
          monthlyTransferIn={monthlyDashboard.monthlyTransferIn}
          monthlyTransferOut={monthlyDashboard.monthlyTransferOut}
          isBalanceHidden={isMainBalanceHidden}
        />
      </DashboardChartsCarousel>
    </DashboardLayout>
  )
}
