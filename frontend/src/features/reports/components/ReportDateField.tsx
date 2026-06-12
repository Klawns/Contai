import { useState } from 'react'
import type { CSSProperties } from 'react'
import { DayPicker } from '@daypicker/react'
import { ptBR } from '@daypicker/react/locale'
import { CalendarDays, ChevronRight } from 'lucide-react'
import { TransactionFieldRow } from '../../transactions/components/FormFields.tsx'
import { SelectionSheet } from '../../transactions/components/SelectionSheet.tsx'
import { fromDateInputValue, toDateInputValue } from '../../transactions/utils/date.ts'
import { formatReportDateLabel } from '../utils/reportFilters.ts'

type ReportDateFieldProps = {
  label: string
  value: string
  accentColor: string
  onChange: (value: string) => void
}

export function ReportDateField({
  label,
  value,
  accentColor,
  onChange,
}: ReportDateFieldProps) {
  const [isOpen, setIsOpen] = useState(false)
  const selectedDate = fromDateInputValue(value)

  return (
    <>
      <TransactionFieldRow
        icon={<CalendarDays className="h-5 w-5" aria-hidden="true" />}
        label={label}
      >
        <button
          type="button"
          className="flex min-h-11 w-full cursor-pointer items-center justify-between gap-2 text-left focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
          onClick={() => setIsOpen(true)}
        >
          <span
            className="min-w-0 max-w-full truncate rounded-full px-3 py-1.5 text-[13px] font-semibold text-white"
            style={{ backgroundColor: accentColor }}
          >
            {formatReportDateLabel(value)}
          </span>
          <ChevronRight className="h-4 w-4 flex-none text-[#9a91a5]" aria-hidden="true" />
        </button>
      </TransactionFieldRow>
      <SelectionSheet title={label} isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div
          className="transaction-date-picker rounded-lg border border-[#eee8f3] bg-white p-2"
          style={{ '--date-accent': accentColor } as CSSProperties}
        >
          <DayPicker
            mode="single"
            locale={ptBR}
            selected={selectedDate}
            defaultMonth={selectedDate}
            onSelect={(date) => {
              if (!date) {
                return
              }

              onChange(toDateInputValue(date))
              setIsOpen(false)
            }}
          />
        </div>
      </SelectionSheet>
    </>
  )
}
