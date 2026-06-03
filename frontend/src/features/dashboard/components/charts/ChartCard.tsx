import { useMemo, useState } from 'react'
import { PieChart as PieChartIcon } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import {
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
} from 'recharts'
import type { ExpenseByCategory } from '../../types/dashboard.ts'
import { formatCurrency } from '../../utils/formatters.ts'

type ChartCardProps = {
  expenses: ExpenseByCategory[]
}

type PieDatum = ExpenseByCategory & {
  percentage: number
}

type PieTooltipProps = {
  active?: boolean
  payload?: Array<{
    payload?: PieDatum
  }>
}

function PieTooltip({ active, payload }: PieTooltipProps) {
  const datum = payload?.[0]?.payload as PieDatum | undefined

  if (!active || !datum) {
    return null
  }

  return (
    <div className="rounded-xl border border-[#ece8f2] bg-white px-3 py-2 shadow-[0_10px_28px_rgba(48,39,61,0.12)]">
      <span className="block text-[12px] font-semibold text-[#241a30]">
        {datum.name}
      </span>
      <span className="text-[11px] font-medium text-[#81798b]">
        {formatCurrency(datum.total)}
      </span>
    </div>
  )
}

function getCategoryIdFromPieDatum(datum: unknown) {
  if (typeof datum !== 'object' || datum === null || !('categoryId' in datum)) {
    return ''
  }

  const categoryId = datum.categoryId

  return typeof categoryId === 'string' ? categoryId : ''
}

export function ChartCard({ expenses }: ChartCardProps) {
  const shouldReduceMotion = useReducedMotion()
  const total = expenses.reduce((sum, expense) => sum + expense.total, 0)
  const data = useMemo(
    () =>
      expenses.map((expense) => ({
        ...expense,
        percentage: total > 0 ? (expense.total / total) * 100 : 0,
      })),
    [expenses, total],
  )
  const [selectedCategoryId, setSelectedCategoryId] = useState(
    data[0]?.categoryId ?? '',
  )
  const selectedCategory =
    data.find((expense) => expense.categoryId === selectedCategoryId) ?? data[0]

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
            Despesas por categoria
          </h3>
          <span className="text-[12px] font-medium text-[#8b8394]">
            {formatCurrency(total)}
          </span>
        </div>
        <span className="grid h-9 w-9 place-items-center rounded-full bg-[#f2eff8] text-[#6a22e5]">
          <PieChartIcon className="h-5 w-5" aria-hidden="true" />
        </span>
      </div>

      <div className="grid gap-3">
        <div className="h-[178px] min-w-0">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart margin={{ top: 0, right: 4, bottom: 0, left: 4 }}>
              <Tooltip content={<PieTooltip />} />
              <Pie
                data={data}
                dataKey="total"
                nameKey="name"
                cx="50%"
                cy="50%"
                innerRadius="58%"
                outerRadius="86%"
                paddingAngle={3}
                cornerRadius={7}
                stroke="none"
                onMouseEnter={(datum) => {
                  const categoryId = getCategoryIdFromPieDatum(datum)

                  if (categoryId) {
                    setSelectedCategoryId(categoryId)
                  }
                }}
                onClick={(datum) => {
                  const categoryId = getCategoryIdFromPieDatum(datum)

                  if (categoryId) {
                    setSelectedCategoryId(categoryId)
                  }
                }}
                isAnimationActive={!shouldReduceMotion}
              >
                {data.map((expense) => (
                  <Cell
                    key={expense.categoryId}
                    fill={expense.color}
                    opacity={
                      selectedCategory?.categoryId === expense.categoryId ? 1 : 0.62
                    }
                    className="cursor-pointer outline-none transition-opacity"
                    tabIndex={0}
                    onFocus={() => setSelectedCategoryId(expense.categoryId)}
                    onClick={() => setSelectedCategoryId(expense.categoryId)}
                  />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="rounded-xl bg-[#fbfafe] px-3 py-3">
          <div className="mt-1 grid min-w-0 grid-cols-[10px_minmax(0,1fr)_auto] items-center gap-2">
            <span
              className="h-2.5 w-2.5 rounded-full"
              style={{ backgroundColor: selectedCategory?.color ?? '#d8d1e3' }}
              aria-hidden="true"
            />
            <strong className="truncate text-[14px] font-semibold text-[#241a30]">
              {selectedCategory?.name ?? 'Sem despesas'}
            </strong>
            <span className="text-[12px] font-semibold text-[#6a22e5]">
              {selectedCategory ? `${selectedCategory.percentage.toFixed(1)}%` : '0%'}
            </span>
          </div>
          <strong className="mt-1 block truncate text-[20px] font-semibold leading-tight text-[#241a30]">
            {formatCurrency(selectedCategory?.total ?? 0)}
          </strong>
        </div>
      </div>
    </motion.article>
  )
}
