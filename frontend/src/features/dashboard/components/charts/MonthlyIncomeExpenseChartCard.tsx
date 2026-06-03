import { BarChart3 } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { MonthlyFinancialSeriesPoint } from '../../types/dashboard.ts'
import { formatCurrency } from '../../utils/formatters.ts'

type MonthlyIncomeExpenseChartCardProps = {
  monthlySeries: MonthlyFinancialSeriesPoint[]
}

function formatCompactCurrency(valueInCents: number) {
  return `R$ ${Math.round(valueInCents / 100000)}k`
}

type ChartTooltipItem = {
  dataKey?: string | number
  name?: string | number
  value?: number | string
}

type IncomeExpenseTooltipProps = {
  active?: boolean
  label?: string | number
  payload?: ChartTooltipItem[]
}

function IncomeExpenseTooltip({
  active,
  label,
  payload,
}: IncomeExpenseTooltipProps) {
  if (!active || !payload?.length) {
    return null
  }

  return (
    <div className="rounded-xl border border-[#ece8f2] bg-white px-3 py-2 shadow-[0_10px_28px_rgba(48,39,61,0.12)]">
      <span className="block text-[12px] font-semibold text-[#241a30]">{label}</span>
      {payload.map((item) => (
        <span
          key={item.dataKey}
          className="block text-[11px] font-medium text-[#81798b]"
        >
          {item.name}: {formatCurrency(Number(item.value ?? 0))}
        </span>
      ))}
    </div>
  )
}

export function MonthlyIncomeExpenseChartCard({
  monthlySeries,
}: MonthlyIncomeExpenseChartCardProps) {
  const shouldReduceMotion = useReducedMotion()

  return (
    <motion.article
      className="min-w-[calc(100vw-32px)] snap-center rounded-[18px] border border-[#ece8f2] bg-white p-4 shadow-[0_16px_38px_rgba(48,39,61,0.07)] sm:min-w-[360px] md:min-w-0"
      whileHover={shouldReduceMotion ? undefined : { y: -2 }}
      whileTap={shouldReduceMotion ? undefined : { scale: 0.99 }}
      transition={{ duration: 0.18, ease: 'easeOut' }}
    >
      <div className="mb-3 flex items-center justify-between gap-3">
        <div>
          <h3 className="m-0 text-[15px] font-semibold leading-tight text-[#241a30]">
            Receitas e despesas
          </h3>
          <span className="text-[12px] font-medium text-[#8b8394]">
            Comparativo mensal
          </span>
        </div>
        <span className="grid h-9 w-9 place-items-center rounded-full bg-[#eef8f3] text-[#17a760]">
          <BarChart3 className="h-5 w-5" aria-hidden="true" />
        </span>
      </div>

      <div className="h-[224px] min-w-0">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={monthlySeries}
            margin={{ top: 8, right: 0, bottom: 0, left: -20 }}
            barCategoryGap={14}
          >
            <CartesianGrid stroke="#f1edf6" vertical={false} />
            <XAxis
              dataKey="monthLabel"
              tickLine={false}
              axisLine={false}
              tick={{ fill: '#8b8394', fontSize: 11, fontWeight: 600 }}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tick={{ fill: '#b1a9bc', fontSize: 10, fontWeight: 600 }}
              tickFormatter={formatCompactCurrency}
              width={44}
            />
            <Tooltip content={<IncomeExpenseTooltip />} cursor={{ fill: '#fbfafe' }} />
            <Bar
              dataKey="income"
              name="Receitas"
              fill="#17a760"
              radius={[6, 6, 0, 0]}
              maxBarSize={18}
              isAnimationActive={!shouldReduceMotion}
            />
            <Bar
              dataKey="expense"
              name="Despesas"
              fill="#e44545"
              radius={[6, 6, 0, 0]}
              maxBarSize={18}
              isAnimationActive={!shouldReduceMotion}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </motion.article>
  )
}
