import { fromDateInputValue } from '../../transactions/utils/date.ts'
import type { RecurrenceFrequency } from '../types/commitments.ts'

function getRecurrenceIntervalCopy(frequency: RecurrenceFrequency, interval: number) {
  const normalizedInterval = Number.isFinite(interval) && interval > 0 ? interval : 1
  const frequencyCopy: Record<
    RecurrenceFrequency,
    { adverb: string; singular: string; plural: string }
  > = {
    daily: { adverb: 'diariamente', singular: 'dia', plural: 'dias' },
    weekly: { adverb: 'semanalmente', singular: 'semana', plural: 'semanas' },
    monthly: { adverb: 'mensalmente', singular: 'mes', plural: 'meses' },
    yearly: { adverb: 'anualmente', singular: 'ano', plural: 'anos' },
  }
  const copy = frequencyCopy[frequency]
  const unit = normalizedInterval === 1 ? copy.singular : copy.plural

  return `${copy.adverb} a cada ${normalizedInterval} ${unit}`
}

function formatRecurrenceEndDate(value: string) {
  if (!value) {
    return ''
  }

  return new Intl.DateTimeFormat('pt-BR').format(fromDateInputValue(value))
}

export function getRecurrenceSummary({
  hasRecurrence,
  frequency,
  interval,
  endsOn,
}: {
  hasRecurrence: boolean
  frequency: RecurrenceFrequency
  interval: number
  endsOn: string
}) {
  if (!hasRecurrence) {
    return 'Este compromisso sera lancado apenas uma vez.'
  }

  const baseSummary = `Repete ${getRecurrenceIntervalCopy(frequency, interval)}`
  const formattedEndDate = formatRecurrenceEndDate(endsOn)

  return formattedEndDate ? `${baseSummary} ate ${formattedEndDate}.` : `${baseSummary}.`
}
