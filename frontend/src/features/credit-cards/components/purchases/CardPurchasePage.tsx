import { useState } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { ReceiptText } from 'lucide-react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { AddCategoryForm } from '../../../transactions/components/AddCategoryForm.tsx'
import {
  DateInput,
  DescriptionInput,
  FormActionButton,
  HeaderAmountInput,
  NoteInput,
  TransactionFieldRow,
} from '../../../transactions/components/FormFields.tsx'
import { CategorySelector } from '../../../transactions/components/Selectors.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import { toDateInputValue } from '../../../transactions/utils/date.ts'
import { useCreateCardPurchase } from '../../hooks/useCardPurchaseMutations.ts'
import { getInstallmentText } from '../../lib/installmentText.ts'
import { toCardPurchasePayload } from '../../lib/formPayloads.ts'
import { purchaseFormSchema } from '../../schemas/purchase.schemas.ts'
import type { PurchaseFormValues } from '../../types/form.types.ts'
import { CardSelector } from '../cards/CardSelector.tsx'
import { ErrorMessage } from '../shared/PageState.tsx'

export function CardPurchasePage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const initialCardId = searchParams.get('cardId') ?? ''
  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<PurchaseFormValues>({
    defaultValues: {
      cardId: initialCardId,
      categoryId: '',
      description: '',
      totalAmount: 0,
      purchaseDate: toDateInputValue(new Date()),
      installmentCount: 1,
      note: '',
    },
    resolver: zodResolver(purchaseFormSchema),
  })
  const selectedCardId = useWatch({ control, name: 'cardId' })
  const selectedAmount = useWatch({ control, name: 'totalAmount' })
  const installmentCount = useWatch({ control, name: 'installmentCount' })
  const mutation = useCreateCardPurchase(selectedCardId)

  return (
    <>
      <TransactionsPageLayout variant="create" tone="expense" animationKey="card-purchase">
        <form
          className="mx-auto min-h-svh w-full max-w-[520px] bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none"
          onSubmit={handleSubmit((values) => {
            mutation.mutate(toCardPurchasePayload(values), { onSuccess: () => navigate('/credit-cards') })
          })}
        >
          <div className="w-full bg-[#c72f4d] px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7">
            <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center">
              <button type="button" className="justify-self-start text-[14px] font-semibold text-white/88" onClick={() => navigate('/credit-cards')}>
                Cancelar
              </button>
              <h1 className="truncate text-center text-[15px] font-semibold">Nova compra</h1>
            </div>
            <div className="mt-10">
              <Controller control={control} name="totalAmount" render={({ field }) => (
                <HeaderAmountInput label="Valor" value={field.value} error={errors.totalAmount?.message} onChange={field.onChange} />
              )} />
            </div>
          </div>
          <div className="-mt-6 overflow-hidden rounded-t-[28px] bg-white">
            <Controller control={control} name="purchaseDate" render={({ field }) => (
              <DateInput label="Compra" value={field.value} accentColor="#c72f4d" error={errors.purchaseDate?.message} onChange={field.onChange} />
            )} />
            <Controller control={control} name="description" render={({ field }) => (
              <DescriptionInput value={field.value} error={errors.description?.message} onChange={field.onChange} />
            )} />
            <Controller control={control} name="cardId" render={({ field }) => (
              <CardSelector value={field.value} error={errors.cardId?.message} onChange={field.onChange} />
            )} />
            <Controller control={control} name="categoryId" render={({ field }) => (
              <CategorySelector type="expense" value={field.value} error={errors.categoryId?.message} onChange={field.onChange} onAddCategory={() => setIsAddingCategory(true)} />
            )} />
            <Controller control={control} name="installmentCount" render={({ field }) => (
              <TransactionFieldRow icon={<ReceiptText className="h-5 w-5" aria-hidden="true" />} label="Parcelas" error={errors.installmentCount?.message}>
                <input
                  type="number"
                  min={1}
                  max={48}
                  className="h-11 w-24 rounded-lg border border-[#e5deee] bg-white px-3 text-[14px] font-semibold text-[#2c2237] outline-none"
                  value={field.value}
                  onChange={(event) => field.onChange(Number(event.target.value))}
                />
                <span className="ml-3 truncate text-[13px] font-semibold text-[#81788c]">
                  {getInstallmentText(selectedAmount, installmentCount)}
                </span>
              </TransactionFieldRow>
            )} />
            <Controller control={control} name="note" render={({ field }) => <NoteInput value={field.value} onChange={field.onChange} />} />
            {mutation.isError ? <ErrorMessage>Nao foi possivel registrar a compra.</ErrorMessage> : null}
            <FormActionButton isPending={mutation.isPending}>Registrar compra</FormActionButton>
          </div>
        </form>
      </TransactionsPageLayout>
      <AddCategoryForm type="expense" isOpen={isAddingCategory} onClose={() => setIsAddingCategory(false)} onCreated={(categoryId) => setValue('categoryId', categoryId, { shouldValidate: true })} />
    </>
  )
}
