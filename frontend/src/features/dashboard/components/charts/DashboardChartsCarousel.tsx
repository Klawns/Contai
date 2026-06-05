import type { ReactNode } from 'react'
import { DashboardEmptyStateCard } from '../cards/DashboardEmptyStateCard.tsx'
import { ChartCard } from './ChartCard.tsx'
import { MonthlyBalanceEvolutionCard } from './MonthlyBalanceEvolutionCard.tsx'
import { MonthlyIncomeExpenseChartCard } from './MonthlyIncomeExpenseChartCard.tsx'
import type {
  ExpenseByCategory,
  MonthlyFinancialSeriesPoint,
} from '../../types/dashboard.ts'

type DashboardChartsCarouselProps = {
  expensesByCategory: ExpenseByCategory[]
  monthlySeries: MonthlyFinancialSeriesPoint[]
  children: ReactNode
}

export function DashboardChartsCarousel({
  expensesByCategory,
  monthlySeries,
  children,
}: DashboardChartsCarouselProps) {
  const hasNoChartData =
    expensesByCategory.length === 0 &&
    monthlySeries.every((point) => point.income === 0 && point.expense === 0)

  return (
    <section className="grid gap-2.5" aria-label="Graficos e resumo mensal">
      <h2 className="m-0 px-1 text-[13px] font-medium leading-none text-[#81798b]">
        Graficos
      </h2>
      <div className="-mx-4 flex snap-x snap-mandatory gap-3 overflow-x-auto px-4 pb-2 sm:-mx-5 sm:px-5 md:mx-0 md:grid md:grid-cols-2 md:gap-4 md:overflow-visible md:px-0 md:pb-0 xl:grid-cols-3">
        {hasNoChartData ? (
          <DashboardEmptyStateCard
            title="Ainda nao ha dados para os graficos"
            message="Cadastre uma transacao para gerar os graficos deste periodo."
            className="min-h-[292px] min-w-[calc(100vw-32px)] snap-center sm:min-w-[360px] md:col-span-2 md:min-w-0 xl:col-span-3"
          />
        ) : (
          <>
            <ChartCard expenses={expensesByCategory} />
            <MonthlyIncomeExpenseChartCard monthlySeries={monthlySeries} />
            <MonthlyBalanceEvolutionCard monthlySeries={monthlySeries} />
            {children}
          </>
        )}
      </div>
    </section>
  )
}
