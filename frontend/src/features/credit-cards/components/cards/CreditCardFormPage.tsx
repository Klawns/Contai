import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  AccountSelector,
} from '../../../transactions/components/Selectors.tsx'
import {
  DescriptionInput,
  FormActionButton,
  HeaderAmountInput,
} from '../../../transactions/components/FormFields.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import { useCreditCards } from '../../hooks/useCreditCards.ts'
import { useCreateCreditCard, useUpdateCreditCard } from '../../hooks/useCreditCardMutations.domain.ts'
import { findById } from '../../lib/creditCardCollections.ts'
import { toCreditCardPayload } from '../../lib/formPayloads.ts'
import { cardFormSchema } from '../../schemas/credit-card.schemas.ts'
import type { CreditCard } from '../../types/credit-card.types.ts'
import type { CardFormValues } from '../../types/form.types.ts'
import { ErrorMessage, FormInvalid, FormLoading } from '../shared/PageState.tsx'
import { NumberRow } from '../shared/NumberRow.tsx'
import { CardStatusSelector } from './CardSelector.tsx'

export function CreditCardFormPage({ mode = 'create' }: { mode?: 'create' | 'edit' }) {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const cardId = searchParams.get('cardId') ?? ''
  const cardsQuery = useCreditCards()
  const card = findById(cardsQuery.data, cardId)

  if (mode === 'edit' && cardsQuery.isLoading) {
    return <FormLoading message="Carregando cartao..." />
  }

  if (mode === 'edit' && (!cardId || !card)) {
    return <FormInvalid message="Cartao nao encontrado." onBack={() => navigate('/credit-cards')} />
  }

  return <CreditCardForm mode={mode} card={card} />
}

function CreditCardForm({ mode, card }: { mode: 'create' | 'edit'; card?: CreditCard }) {
  const navigate = useNavigate()
  const createMutation = useCreateCreditCard()
  const updateMutation = useUpdateCreditCard(card?.id ?? '')
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<CardFormValues>({
    defaultValues: card
      ? {
          name: card.name,
          linkedAccountId: card.linkedAccountId,
          limitTotal: card.limitTotal,
          closingDay: card.closingDay,
          dueDay: card.dueDay,
          status: card.status,
        }
      : {
          name: '',
          linkedAccountId: '',
          limitTotal: 0,
          closingDay: 10,
          dueDay: 15,
          status: 'active',
        },
    resolver: zodResolver(cardFormSchema),
  })
  const isPending = mode === 'edit' ? updateMutation.isPending : createMutation.isPending
  const hasError = mode === 'edit' ? updateMutation.isError : createMutation.isError

  return (
    <TransactionsPageLayout variant="create" tone="expense" animationKey={`card-${mode}-${card?.id ?? 'new'}`}>
      <form
        className="scrollbar-none mx-auto h-full min-h-0 w-full max-w-[520px] overflow-y-auto overflow-x-hidden bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none"
        onSubmit={handleSubmit((values) => {
          const payload = toCreditCardPayload(values, mode)
          const onSuccess = () => navigate('/credit-cards')

          if (mode === 'edit') {
            updateMutation.mutate(payload, { onSuccess })
            return
          }
          createMutation.mutate(payload, { onSuccess })
        })}
      >
        <div className="w-full bg-[#216fb8] px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7">
          <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center">
            <button type="button" className="justify-self-start text-[14px] font-semibold text-white/88" onClick={() => navigate('/credit-cards')}>
              Cancelar
            </button>
            <h1 className="truncate text-center text-[15px] font-semibold">{mode === 'edit' ? 'Editar cartao' : 'Novo cartao'}</h1>
          </div>
          <div className="mt-10">
            <Controller control={control} name="limitTotal" render={({ field }) => (
              <HeaderAmountInput label="Limite" value={field.value} error={errors.limitTotal?.message} onChange={field.onChange} />
            )} />
          </div>
        </div>
        <div className="-mt-6 rounded-t-[28px] bg-white">
          <Controller control={control} name="name" render={({ field }) => (
            <DescriptionInput value={field.value} error={errors.name?.message} onChange={field.onChange} />
          )} />
          <Controller control={control} name="linkedAccountId" render={({ field }) => (
            <AccountSelector label="Conta" value={field.value} error={errors.linkedAccountId?.message} onChange={field.onChange} />
          )} />
          <Controller control={control} name="closingDay" render={({ field }) => (
            <NumberRow label="Fechamento" value={field.value} error={errors.closingDay?.message} onChange={field.onChange} />
          )} />
          <Controller control={control} name="dueDay" render={({ field }) => (
            <NumberRow label="Vencimento" value={field.value} error={errors.dueDay?.message} onChange={field.onChange} />
          )} />
          {mode === 'edit' ? (
            <Controller control={control} name="status" render={({ field }) => (
              <CardStatusSelector
                value={field.value}
                error={errors.status?.message}
                onChange={field.onChange}
              />
            )} />
          ) : null}
          {hasError ? <ErrorMessage>Nao foi possivel salvar o cartao.</ErrorMessage> : null}
          <FormActionButton isPending={isPending}>{mode === 'edit' ? 'Salvar alteracoes' : 'Salvar cartao'}</FormActionButton>
        </div>
      </form>
    </TransactionsPageLayout>
  )
}
