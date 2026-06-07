import { zodResolver } from '@hookform/resolvers/zod'
import { useMemo, useState } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { AddCategoryForm } from '../../../transactions/components/AddCategoryForm.tsx'
import {
  DateInput,
  DescriptionInput,
  FormActionButton,
  HeaderAmountInput,
  NoteInput,
} from '../../../transactions/components/FormFields.tsx'
import {
  AccountSelector,
  CategorySelector,
} from '../../../transactions/components/Selectors.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import {
  useCreateCommitment,
  useUpdateCommitment,
} from '../../hooks/useCommitmentMutations.ts'
import {
  getDefaultFormValues,
  getInitialFormValues,
  toCommitmentPayload,
} from '../../lib/commitmentFormValues.ts'
import { typeCopy } from '../../lib/commitmentPresentation.ts'
import { categoryTypeForCommitment } from '../../lib/commitmentType.ts'
import { commitmentFormSchema } from '../../schemas/commitmentForm.ts'
import type { CommitmentFormValues } from '../../types/commitmentForm.ts'
import type { Commitment, CommitmentType } from '../../types/commitments.ts'
import { RecurrenceSection } from './RecurrenceSection.tsx'

export function CommitmentForm({
  type,
  mode,
  initialCommitment,
}: {
  type: CommitmentType
  mode: 'create' | 'edit'
  initialCommitment?: Commitment
}) {
  const navigate = useNavigate()
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const createMutation = useCreateCommitment(type)
  const updateMutation = useUpdateCommitment(type, initialCommitment?.id ?? '')
  const categoryType = categoryTypeForCommitment(type)
  const styles = typeCopy[type]
  const defaultValues = useMemo(
    () =>
      mode === 'edit' && initialCommitment
        ? getInitialFormValues(initialCommitment)
        : getDefaultFormValues(),
    [initialCommitment, mode],
  )
  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<CommitmentFormValues>({
    defaultValues,
    resolver: zodResolver(commitmentFormSchema),
  })
  const hasRecurrence = useWatch({ control, name: 'hasRecurrence' })
  const isPending = mode === 'edit' ? updateMutation.isPending : createMutation.isPending
  const hasMutationError = mode === 'edit' ? updateMutation.isError : createMutation.isError

  return (
    <>
      <TransactionsPageLayout
        variant="create"
        tone={styles.tone}
        animationKey={`${type}-${mode}-${initialCommitment?.id ?? 'new'}`}
      >
        <form
          className="mx-auto min-h-svh w-full max-w-[520px] bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none"
          onSubmit={handleSubmit((values) => {
            const payload = toCommitmentPayload(values)
            const onSuccess = () => navigate('/planning')

            if (mode === 'edit') {
              updateMutation.mutate(payload, { onSuccess })
              return
            }
            createMutation.mutate(payload, { onSuccess })
          })}
        >
          <div className={`${styles.header} w-full px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7 lg:px-10 lg:pb-14`}>
            <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center">
              <button
                type="button"
                className="justify-self-start text-[14px] font-semibold text-white/88 transition-colors hover:text-white focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-white"
                onClick={() => navigate('/planning')}
              >
                Cancelar
              </button>
              <h1 className="truncate text-center text-[15px] font-semibold">
                {mode === 'edit' ? styles.editTitle : styles.newTitle}
              </h1>
            </div>
            <div className="mt-10">
              <Controller
                control={control}
                name="amount"
                render={({ field }) => (
                  <HeaderAmountInput
                    label="Valor"
                    value={field.value}
                    error={errors.amount?.message}
                    onChange={field.onChange}
                  />
                )}
              />
            </div>
          </div>

          <div className="-mt-6 overflow-hidden rounded-t-[28px] bg-white md:rounded-t-[32px]">
            <Controller
              control={control}
              name="dueOn"
              render={({ field }) => (
                <DateInput
                  label="Vencimento"
                  value={field.value}
                  accentColor={styles.accent}
                  error={errors.dueOn?.message}
                  onChange={field.onChange}
                />
              )}
            />
            <Controller
              control={control}
              name="description"
              render={({ field }) => (
                <DescriptionInput
                  value={field.value}
                  error={errors.description?.message}
                  onChange={field.onChange}
                />
              )}
            />
            <Controller
              control={control}
              name="categoryId"
              render={({ field }) => (
                <CategorySelector
                  type={categoryType}
                  value={field.value}
                  error={errors.categoryId?.message}
                  onChange={field.onChange}
                  onAddCategory={() => setIsAddingCategory(true)}
                />
              )}
            />
            <Controller
              control={control}
              name="accountId"
              render={({ field }) => (
                <AccountSelector
                  label="Conta"
                  value={field.value}
                  error={errors.accountId?.message}
                  onChange={field.onChange}
                />
              )}
            />
            <Controller
              control={control}
              name="note"
              render={({ field }) => <NoteInput value={field.value} onChange={field.onChange} />}
            />
            <RecurrenceSection
              control={control}
              errors={errors}
              hasRecurrence={hasRecurrence}
              accentColor={styles.accent}
            />
            {hasMutationError ? (
              <p className="mx-4 mt-4 rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838] md:mx-5">
                Nao foi possivel salvar o compromisso.
              </p>
            ) : null}
            <FormActionButton isPending={isPending}>
              {mode === 'edit' ? 'Salvar alteracoes' : 'Salvar compromisso'}
            </FormActionButton>
          </div>
        </form>
      </TransactionsPageLayout>

      <AddCategoryForm
        type={categoryType}
        isOpen={isAddingCategory}
        onClose={() => setIsAddingCategory(false)}
        onCreated={(categoryId) => setValue('categoryId', categoryId, { shouldValidate: true })}
      />
    </>
  )
}
