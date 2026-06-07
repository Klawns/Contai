import { zodResolver } from '@hookform/resolvers/zod'
import { PencilLine } from 'lucide-react'
import { useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { AddCategoryForm } from '../../../transactions/components/AddCategoryForm.tsx'
import {
  DateInput,
  FormActionButton,
  HeaderAmountInput,
  NoteInput,
  TransactionFieldRow,
} from '../../../transactions/components/FormFields.tsx'
import {
  AccountSelector,
  CategorySelector,
} from '../../../transactions/components/Selectors.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import { useSettleCommitment } from '../../hooks/useCommitmentMutations.ts'
import { typeCopy } from '../../lib/commitmentPresentation.ts'
import { categoryTypeForCommitment } from '../../lib/commitmentType.ts'
import {
  getSettlementDefaultValues,
  toSettlementPayload,
} from '../../lib/settlementFormValues.ts'
import { settlementFormSchema } from '../../schemas/settlementForm.ts'
import type { Commitment, CommitmentType } from '../../types/commitments.ts'
import type { SettlementFormValues } from '../../types/settlementForm.ts'

export function SettlementForm({
  type,
  commitment,
}: {
  type: CommitmentType
  commitment: Commitment
}) {
  const navigate = useNavigate()
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const settleMutation = useSettleCommitment(type, commitment.id)
  const categoryType = categoryTypeForCommitment(type)
  const styles = typeCopy[type]
  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<SettlementFormValues>({
    defaultValues: getSettlementDefaultValues(commitment),
    resolver: zodResolver(settlementFormSchema),
  })

  return (
    <>
      <TransactionsPageLayout variant="create" tone={styles.tone} animationKey={`settle-${commitment.id}`}>
        <form
          className="mx-auto min-h-svh w-full max-w-[520px] bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none"
          onSubmit={handleSubmit((values) => {
            settleMutation.mutate(toSettlementPayload(values), {
              onSuccess: () => navigate('/planning'),
            })
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
                {styles.settleTitle}
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
              name="occurredOn"
              render={({ field }) => (
                <DateInput
                  label="Data"
                  value={field.value}
                  accentColor={styles.accent}
                  error={errors.occurredOn?.message}
                  onChange={field.onChange}
                />
              )}
            />
            <TransactionFieldRow
              icon={<PencilLine className="h-5 w-5" aria-hidden="true" />}
              label="Descricao"
            >
              <span className="truncate text-[15px] font-semibold text-[#2c2237]">
                {commitment.description}
              </span>
            </TransactionFieldRow>
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
            {settleMutation.isError ? (
              <p className="mx-4 mt-4 rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838] md:mx-5">
                Nao foi possivel quitar o compromisso.
              </p>
            ) : null}
            <FormActionButton isPending={settleMutation.isPending}>
              {styles.settleButton}
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
