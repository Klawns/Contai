import { useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Landmark } from 'lucide-react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { AddCategoryForm } from '../../../transactions/components/AddCategoryForm.tsx'
import {
  DateInput,
  FormActionButton,
  HeaderAmountInput,
  NoteInput,
  TransactionFieldRow,
} from '../../../transactions/components/FormFields.tsx'
import { CategorySelector } from '../../../transactions/components/Selectors.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import { toDateInputValue } from '../../../transactions/utils/date.ts'
import { useCardInvoice } from '../../hooks/useCardInvoices.ts'
import { usePayCardInvoice } from '../../hooks/useCardInvoiceMutations.ts'
import { toPayInvoicePayload } from '../../lib/formPayloads.ts'
import { payInvoiceFormSchema } from '../../schemas/invoice.schemas.ts'
import type { PayInvoiceFormValues } from '../../types/form.types.ts'
import { ErrorMessage, FormInvalid, FormLoading } from '../shared/PageState.tsx'

export function PayCardInvoicePage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const invoiceId = searchParams.get('invoiceId') ?? ''
  const invoiceQuery = useCardInvoice(invoiceId)
  const payMutation = usePayCardInvoice(invoiceId)
  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<PayInvoiceFormValues>({
    defaultValues: {
      occurredOn: toDateInputValue(new Date()),
      categoryId: '',
      note: '',
    },
    resolver: zodResolver(payInvoiceFormSchema),
  })
  const invoice = invoiceQuery.data

  if (invoiceQuery.isLoading) {
    return <FormLoading message="Carregando fatura..." />
  }

  if (!invoiceId || !invoice || (invoice.effectiveStatus !== 'closed' && invoice.effectiveStatus !== 'overdue')) {
    return <FormInvalid message="Esta fatura nao pode ser paga." onBack={() => navigate('/credit-cards')} />
  }

  return (
    <>
      <TransactionsPageLayout variant="create" tone="expense" animationKey={`pay-invoice-${invoice.id}`}>
        <form
          className="scrollbar-none mx-auto h-full min-h-0 w-full max-w-[520px] overflow-y-auto overflow-x-hidden bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none"
          onSubmit={handleSubmit((values) => {
            payMutation.mutate(toPayInvoicePayload(values), { onSuccess: () => navigate('/transactions') })
          })}
        >
          <div className="w-full bg-[#147a46] px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7">
            <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center">
              <button type="button" className="justify-self-start text-[14px] font-semibold text-white/88" onClick={() => navigate(`/credit-cards/invoice?invoiceId=${encodeURIComponent(invoice.id)}`)}>
                Cancelar
              </button>
              <h1 className="truncate text-center text-[15px] font-semibold">Pagar fatura</h1>
            </div>
            <div className="mt-10">
              <HeaderAmountInput label="Valor" value={invoice.amount} onChange={() => undefined} />
            </div>
          </div>
          <div className="-mt-6 rounded-t-[28px] bg-white">
            <Controller control={control} name="occurredOn" render={({ field }) => (
              <DateInput label="Pagamento" value={field.value} accentColor="#147a46" error={errors.occurredOn?.message} onChange={field.onChange} />
            )} />
            <TransactionFieldRow icon={<Landmark className="h-5 w-5" aria-hidden="true" />} label="Conta">
              <span className="truncate text-[15px] font-semibold text-[#2c2237]">Conta vinculada ao cartao</span>
            </TransactionFieldRow>
            <Controller control={control} name="categoryId" render={({ field }) => (
              <CategorySelector type="expense" value={field.value} error={errors.categoryId?.message} onChange={field.onChange} onAddCategory={() => setIsAddingCategory(true)} />
            )} />
            <Controller control={control} name="note" render={({ field }) => <NoteInput value={field.value} onChange={field.onChange} />} />
            {payMutation.isError ? <ErrorMessage>Nao foi possivel pagar a fatura.</ErrorMessage> : null}
            <FormActionButton isPending={payMutation.isPending}>Pagar fatura</FormActionButton>
          </div>
        </form>
      </TransactionsPageLayout>
      <AddCategoryForm type="expense" isOpen={isAddingCategory} onClose={() => setIsAddingCategory(false)} onCreated={(categoryId) => setValue('categoryId', categoryId, { shouldValidate: true })} />
    </>
  )
}
