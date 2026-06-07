import { formatCurrency } from '../../../transactions/utils/money.ts'

export function PlanningTotals({ totals }: { totals: { open: number; overdue: number } }) {
  return (
    <section className="grid grid-cols-2 border-b border-[#f0ebf6] pb-4">
      <div className="min-w-0 pr-4">
        <span className="block text-[12px] font-semibold leading-tight text-[#81788c]">
          Em aberto
        </span>
        <strong className="mt-1 block truncate text-[18px] font-semibold leading-tight text-[#2c2237]">
          {formatCurrency(totals.open)}
        </strong>
      </div>
      <div className="min-w-0 border-l border-[#f0ebf6] pl-4">
        <span className="block text-[12px] font-semibold leading-tight text-[#81788c]">
          Vencido
        </span>
        <strong className="mt-1 block truncate text-[18px] font-semibold leading-tight text-[#c72f4d]">
          {formatCurrency(totals.overdue)}
        </strong>
      </div>
    </section>
  )
}
