import { useState } from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { DayPicker } from '@daypicker/react'
import { ptBR } from '@daypicker/react/locale'
import {
  CalendarDays,
  FileText,
  MessageSquareText,
  PencilLine,
} from 'lucide-react'
import { formatCentsAsInput, parseCurrencyToCents } from '../utils/money.ts'
import { fromDateInputValue, toDateInputValue } from '../utils/date.ts'
import { SelectionSheet } from './SelectionSheet.tsx'

type FieldProps = {
  label: string
  error?: string
  children: ReactNode
}

export function FormSection({ label, error, children }: FieldProps) {
  return (
    <label className="grid gap-2">
      <span className="text-[13px] font-semibold text-[#5f536d]">{label}</span>
      {children}
      {error ? <span className="text-[12px] font-medium text-[#c72f4d]">{error}</span> : null}
    </label>
  )
}

const inputClasses =
  'h-12 w-full rounded-lg border border-[#e5deee] bg-white px-4 text-[15px] font-medium text-[#2c2237] outline-none transition-colors placeholder:text-[#aaa2b4] focus:border-[#7b2cff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]'

type AmountInputProps = {
  value: number
  onChange: (value: number) => void
  error?: string
}

type HeaderAmountInputProps = AmountInputProps & {
  label: string
}

export function AmountInput({ value, onChange, error }: AmountInputProps) {
  return (
    <FormSection label="Valor" error={error}>
      <input
        inputMode="numeric"
        className={`${inputClasses} text-[20px] font-semibold`}
        value={formatCentsAsInput(value)}
        onChange={(event) => onChange(parseCurrencyToCents(event.target.value))}
      />
    </FormSection>
  )
}

export function HeaderAmountInput({ label, value, onChange, error }: HeaderAmountInputProps) {
  return (
    <div className="grid gap-2 text-white">
      <label className="grid cursor-text gap-2 border-b border-white/24 pb-2 transition-colors focus-within:border-white/60">
        <span className="text-left text-[14px] font-semibold text-white/78">{label}</span>
        <input
          inputMode="numeric"
          aria-label={label}
          className="w-full cursor-text bg-transparent text-left text-[76px] font-bold leading-none text-white outline-none placeholder:text-white/55 sm:text-[92px] md:text-[104px] lg:text-[112px]"
          value={formatCentsAsInput(value)}
          placeholder="0,00"
          onChange={(event) => onChange(parseCurrencyToCents(event.target.value))}
        />
      </label>
      {error ? (
        <span className="text-left text-[12px] font-semibold text-white">{error}</span>
      ) : null}
    </div>
  )
}

type DateInputProps = {
  value: string
  accentColor: string
  label?: string
  onChange: (value: string) => void
  error?: string
}

export function DateInput({ value, accentColor, label = 'Data', onChange, error }: DateInputProps) {
  const [isCalendarOpen, setIsCalendarOpen] = useState(false)
  const today = toDateInputValue(new Date())
  const yesterdayDate = new Date()
  yesterdayDate.setDate(yesterdayDate.getDate() - 1)
  const yesterday = toDateInputValue(yesterdayDate)
  const isOtherSelected = Boolean(value && value !== today && value !== yesterday)
  const selectedDate = value ? fromDateInputValue(value) : undefined

  return (
    <>
      <TransactionFieldRow
        icon={<CalendarDays className="h-5 w-5" aria-hidden="true" />}
        label={label}
        error={error}
      >
        <div className="flex min-h-11 flex-wrap items-center justify-start gap-2">
          <DateChip
            accentColor={accentColor}
            isSelected={value === today}
            onClick={() => onChange(today)}
          >
            Hoje
          </DateChip>
          <DateChip
            accentColor={accentColor}
            isSelected={value === yesterday}
            onClick={() => onChange(yesterday)}
          >
            Ontem
          </DateChip>
          <DateChip
            accentColor={accentColor}
            isSelected={isOtherSelected}
            onClick={() => setIsCalendarOpen(true)}
          >
            Outros
          </DateChip>
        </div>
      </TransactionFieldRow>
      <SelectionSheet title="Selecionar data" isOpen={isCalendarOpen} onClose={() => setIsCalendarOpen(false)}>
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
              setIsCalendarOpen(false)
            }}
          />
        </div>
      </SelectionSheet>
    </>
  )
}

function DateChip({
  isSelected,
  accentColor,
  onClick,
  children,
}: {
  isSelected: boolean
  accentColor: string
  onClick: () => void
  children: ReactNode
}) {
  return (
    <button
      type="button"
      className={`rounded-full px-3 py-1.5 text-[12px] font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
        isSelected ? 'text-white' : 'bg-[#f4f1f7] text-[#5f536d] hover:bg-[#ece6f3]'
      }`}
      style={isSelected ? { backgroundColor: accentColor } : undefined}
      onClick={onClick}
    >
      {children}
    </button>
  )
}

type DescriptionInputProps = {
  value: string
  onChange: (value: string) => void
  error?: string
}

export function DescriptionInput({
  value,
  onChange,
  error,
}: DescriptionInputProps) {
  return (
    <TransactionFieldRow
      icon={<PencilLine className="h-5 w-5" aria-hidden="true" />}
      label="Descricao"
      error={error}
    >
      <input
        className="h-11 w-full bg-transparent text-left text-[15px] font-semibold text-[#2c2237] outline-none placeholder:text-[#aaa2b4]"
        value={value}
        placeholder="Adicionar descricao"
        onChange={(event) => onChange(event.target.value)}
      />
    </TransactionFieldRow>
  )
}

type NoteInputProps = {
  value: string
  onChange: (value: string) => void
}

export function NoteInput({ value, onChange }: NoteInputProps) {
  return (
    <TransactionFieldRow
      icon={<MessageSquareText className="h-5 w-5" aria-hidden="true" />}
      label="Observacao"
    >
      <input
        className="h-11 w-full bg-transparent text-left text-[15px] font-semibold text-[#2c2237] outline-none placeholder:text-[#aaa2b4]"
        value={value}
        placeholder="Opcional"
        onChange={(event) => onChange(event.target.value)}
      />
    </TransactionFieldRow>
  )
}

type FormActionButtonProps = {
  isPending: boolean
  children: ReactNode
}

export function FormActionButton({ isPending, children }: FormActionButtonProps) {
  return (
    <div className="sticky bottom-[var(--app-mobile-sticky-bottom)] bg-white/96 px-4 pb-4 pt-5 backdrop-blur md:static md:bg-transparent md:px-5 md:pb-5 md:pt-6 md:backdrop-blur-none">
      <button
        type="submit"
        disabled={isPending}
        className="mx-auto block h-12 w-full max-w-[420px] cursor-pointer rounded-lg bg-[#281d35] px-4 text-[15px] font-semibold text-white shadow-[0_6px_14px_rgba(40,29,53,0.10)] transition-colors hover:bg-[#3a2a4a] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] disabled:cursor-not-allowed disabled:opacity-65"
      >
        {isPending ? 'Salvando...' : children}
      </button>
    </div>
  )
}

type TransactionFieldRowProps = {
  icon?: ReactNode
  label: string
  error?: string
  className?: string
  children: ReactNode
}

export function TransactionFieldRow({
  icon,
  label,
  error,
  className = '',
  children,
}: TransactionFieldRowProps) {
  return (
    <div className={`border-b border-[#f0ebf5] px-4 py-2 md:px-5 ${className}`}>
      <div className="flex min-h-[56px] items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-full bg-[#f4f1f7] text-[#6f647b]">
          {icon ?? <FileText className="h-5 w-5" aria-hidden="true" />}
        </span>
        <span className="min-w-[88px] flex-none text-[14px] font-semibold text-[#5f536d] md:min-w-[96px]">
          {label}
        </span>
        <div className="flex min-w-0 flex-1 items-center">{children}</div>
      </div>
      {error ? (
        <p className="pb-2 pl-[48px] text-[12px] font-semibold text-[#c72f4d] md:pl-[52px]">
          {error}
        </p>
      ) : null}
    </div>
  )
}
