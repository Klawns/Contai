import { TrendingUp } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { MonthlyFinancialSeriesPoint } from '../../types/dashboard.ts'
import { getBalanceTone } from '../../utils/balance.ts'
import { formatCurrency } from '../../utils/formatters.ts'

type MonthlyBalanceEvolutionCardProps = {
  monthlySeries: MonthlyFinancialSeriesPoint[]
}

function formatCompactCurrency(valueInCents: number) {
  const signal = valueInCents < 0 ? '-' : ''

  return `${signal}R$ ${Math.round(Math.abs(valueInCents) / 100000)}k`
}

type BalanceTooltipProps = {
  active?: boolean
  label?: string | number
  payload?: Array<{
    value?: number | string
  }>
}

function BalanceTooltip({ active, label, payload }: BalanceTooltipProps) {
  const value = Number(payload?.[0]?.value ?? 0)

  if (!active || !payload?.length) {
    return null
  }

  return (
    <div className="rounded-xl border border-[#ece8f2] bg-white px-3 py-2 shadow-[0_10px_28px_rgba(48,39,61,0.12)]">
      <span className="block text-[12px] font-semibold text-[#241a30]">{label}</span>
      <span className="text-[11px] font-medium text-[#81798b]">
        Saldo: {formatCurrency(value)}
      </span>
    </div>
  )
}

export function MonthlyBalanceEvolutionCard({
  monthlySeries,
}: MonthlyBalanceEvolutionCardProps) {
  const shouldReduceMotion = useReducedMotion()
  const latestBalance = monthlySeries.at(-1)?.balance ?? 0
  const balanceTone = getBalanceTone(latestBalance)

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
            Evolucao do saldo
          </h3>
          <span className={`text-[12px] font-semibold ${balanceTone.textClass}`}>
            {formatCurrency(latestBalance)}
          </span>
        </div>
        <span className="grid h-9 w-9 place-items-center rounded-full bg-[#f2eff8] text-[#6a22e5]">
          <TrendingUp className="h-5 w-5" aria-hidden="true" />
        </span>
      </div>

      <div className="h-[224px] min-w-0">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart
            data={monthlySeries}
            margin={{ top: 10, right: 4, bottom: 0, left: -20 }}
          >
            <defs>
              <linearGradient id="monthly-balance-fill" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#6a22e5" stopOpacity={0.24} />
                <stop offset="95%" stopColor="#6a22e5" stopOpacity={0.02} />
              </linearGradient>
            </defs>
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
            <Tooltip content={<BalanceTooltip />} cursor={{ stroke: '#d8d1e3' }} />
            <Area
              type="monotone"
              dataKey="balance"
              name="Saldo"
              stroke="#6a22e5"
              strokeWidth={3}
              fill="url(#monthly-balance-fill)"
              dot={{ r: 3, fill: '#6a22e5', strokeWidth: 0 }}
              activeDot={{ r: 5, fill: '#6a22e5', stroke: '#ffffff', strokeWidth: 2 }}
              isAnimationActive={!shouldReduceMotion}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </motion.article>
  )
}
