import { CalendarDays } from 'lucide-react'
import { TransactionFieldRow } from '../../../transactions/components/FormFields.tsx'

export function NumberRow({
  label,
  value,
  error,
  onChange,
}: {
  label: string
  value: number
  error?: string
  onChange: (value: number) => void
}) {
  return (
    <TransactionFieldRow icon={<CalendarDays className="h-5 w-5" aria-hidden="true" />} label={label} error={error}>
      <input
        type="number"
        min={1}
        max={31}
        className="h-11 w-24 rounded-lg border border-[#e5deee] bg-white px-3 text-[14px] font-semibold text-[#2c2237] outline-none"
        value={value}
        onChange={(event) => onChange(Number(event.target.value))}
      />
    </TransactionFieldRow>
  )
}
