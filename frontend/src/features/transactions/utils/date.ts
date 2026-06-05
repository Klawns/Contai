import type { SelectedMonth } from '../../../components/MonthSelector.tsx'

export function formatLocalRFC3339(date: Date) {
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

export function toDateInputValue(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

export function fromDateInputValue(value: string) {
  const [year, month, day] = value.split('-').map(Number)

  return new Date(year, month - 1, day, 12, 0, 0)
}

export function getCurrentSelectedMonth(): SelectedMonth {
  const now = new Date()

  return {
    year: now.getFullYear(),
    monthIndex: now.getMonth(),
  }
}

export function formatMonthQuery(month: SelectedMonth) {
  return `${month.year}-${String(month.monthIndex + 1).padStart(2, '0')}`
}

export function getMonthPeriod(month: SelectedMonth) {
  const startAt = new Date(month.year, month.monthIndex, 1, 0, 0, 0)
  const endAt = new Date(month.year, month.monthIndex + 1, 0, 23, 59, 59)

  return {
    startAt: formatLocalRFC3339(startAt),
    endAt: formatLocalRFC3339(endAt),
  }
}
