import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { Clock, Repeat2 } from 'lucide-react'
import { Controller, useWatch } from 'react-hook-form'
import type { Control, FieldErrors } from 'react-hook-form'
import { TransactionSelectField } from '../../../transactions/components/Selectors.tsx'
import { recurrenceFrequencyOptions } from '../../lib/commitmentPresentation.ts'
import { getRecurrenceSummary } from '../../lib/recurrenceSummary.ts'
import type { CommitmentFormValues } from '../../types/commitmentForm.ts'

type RecurrenceSectionProps = {
  control: Control<CommitmentFormValues>
  errors: FieldErrors<CommitmentFormValues>
  hasRecurrence: boolean
  accentColor: string
}

export function RecurrenceSection({
  control,
  errors,
  hasRecurrence,
  accentColor,
}: RecurrenceSectionProps) {
  const shouldReduceMotion = useReducedMotion()
  const recurrenceFrequency = useWatch({ control, name: 'recurrenceFrequency' }) ?? 'monthly'
  const recurrenceInterval = useWatch({ control, name: 'recurrenceInterval' }) ?? 1
  const recurrenceEndsOn = useWatch({ control, name: 'recurrenceEndsOn' }) ?? ''
  const summary = getRecurrenceSummary({
    hasRecurrence,
    frequency: recurrenceFrequency,
    interval: recurrenceInterval,
    endsOn: recurrenceEndsOn,
  })

  return (
    <section className="border-b border-[#f0ebf5] px-4 py-4 md:px-5">
      <div className="rounded-xl border border-[#eee8f3] bg-[#fbf9fe] px-4 py-3 shadow-[0_10px_28px_rgba(43,35,54,0.06)]">
        <div className="flex items-start gap-3">
          <span
            className="grid h-10 w-10 flex-none place-items-center rounded-full bg-white shadow-[0_6px_18px_rgba(43,35,54,0.08)]"
            style={{ color: accentColor }}
          >
            <Repeat2 className="h-5 w-5" aria-hidden="true" />
          </span>
          <div className="min-w-0 flex-1">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <h2 className="text-[15px] font-semibold text-[#2c2237]">Recorrencia</h2>
              </div>
              <Controller
                control={control}
                name="hasRecurrence"
                render={({ field }) => (
                  <motion.button
                    type="button"
                    role="switch"
                    aria-checked={field.value}
                    aria-label="Ativar recorrencia"
                    className={`relative h-8 w-[52px] flex-none cursor-pointer rounded-full transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
                      field.value ? '' : 'bg-[#d8d3df]'
                    }`}
                    style={field.value ? { backgroundColor: accentColor } : undefined}
                    whileTap={shouldReduceMotion ? undefined : { scale: 0.96 }}
                    onClick={() => field.onChange(!field.value)}
                  >
                    <motion.span
                      className="absolute left-1 top-1 h-6 w-6 rounded-full bg-white shadow-[0_2px_6px_rgba(43,35,54,0.22)]"
                      animate={{ x: field.value ? 20 : 0 }}
                      transition={shouldReduceMotion ? { duration: 0 } : { duration: 0.18 }}
                      aria-hidden="true"
                    />
                  </motion.button>
                )}
              />
            </div>
            {hasRecurrence ? (
              <p className="mt-3 rounded-lg bg-white px-3 py-2 text-[12px] font-semibold leading-relaxed text-[#2c2237]">
                {summary}
              </p>
            ) : null}
          </div>
        </div>

        <AnimatePresence initial={false}>
          {hasRecurrence ? (
            <motion.div
              className="overflow-hidden"
              initial={shouldReduceMotion ? false : { height: 0, opacity: 0, y: -6 }}
              animate={{ height: 'auto', opacity: 1, y: 0 }}
              exit={
                shouldReduceMotion
                  ? { height: 0, opacity: 0, y: 0 }
                  : { height: 0, opacity: 0, y: -6 }
              }
              transition={
                shouldReduceMotion ? { duration: 0 } : { duration: 0.22, ease: 'easeOut' }
              }
            >
              <div className="mt-4 overflow-hidden rounded-lg border border-[#eee8f3] bg-white">
                <Controller
                  control={control}
                  name="recurrenceFrequency"
                  render={({ field }) => (
                    <TransactionSelectField
                      label="Frequencia"
                      value={field.value}
                      placeholder="Selecione"
                      icon={<Clock className="h-5 w-5" aria-hidden="true" />}
                      options={recurrenceFrequencyOptions}
                      onChange={field.onChange}
                      chipClassName="bg-[#f4f1f7] text-[#2c2237]"
                      sheetTitle="Frequencia"
                    />
                  )}
                />
                <div className="grid gap-3 px-4 py-4 md:grid-cols-2">
                  <Controller
                    control={control}
                    name="recurrenceInterval"
                    render={({ field }) => (
                      <label className="grid gap-1">
                        <span className="text-[12px] font-semibold text-[#5f536d]">A cada</span>
                        <input
                          type="number"
                          min={1}
                          className="h-11 rounded-lg border border-[#e5deee] bg-white px-3 text-[14px] font-semibold text-[#2c2237] outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                          value={field.value}
                          onChange={(event) => field.onChange(Number(event.target.value))}
                        />
                        {errors.recurrenceInterval?.message ? (
                          <span className="text-[12px] font-medium text-[#c72f4d]">
                            {errors.recurrenceInterval.message}
                          </span>
                        ) : null}
                      </label>
                    )}
                  />
                  <Controller
                    control={control}
                    name="recurrenceEndsOn"
                    render={({ field }) => (
                      <label className="grid gap-1">
                        <span className="text-[12px] font-semibold text-[#5f536d]">Encerrar em</span>
                        <input
                          type="date"
                          className="h-11 rounded-lg border border-[#e5deee] bg-white px-3 text-[14px] font-semibold text-[#2c2237] outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                          value={field.value}
                          onChange={(event) => field.onChange(event.target.value)}
                        />
                        {errors.recurrenceEndsOn?.message ? (
                          <span className="text-[12px] font-medium text-[#c72f4d]">
                            {errors.recurrenceEndsOn.message}
                          </span>
                        ) : null}
                      </label>
                    )}
                  />
                </div>
              </div>
            </motion.div>
          ) : null}
        </AnimatePresence>
      </div>
    </section>
  )
}
