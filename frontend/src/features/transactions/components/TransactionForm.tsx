import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useCreateTransaction } from '../hooks/useCreateTransaction.ts'
import { useUpdateTransaction } from '../hooks/useUpdateTransaction.ts'
import type {
  CategoryTransactionType,
  CreateExpenseTransactionPayload,
  CreateIncomeTransactionPayload,
  CreateTransferTransactionPayload,
  Transaction,
  TransactionType,
  UpdateTransactionPayloadByType,
} from '../types/transactions.ts'
import { fromDateInputValue, formatLocalRFC3339, toDateInputValue } from '../utils/date.ts'
import {
  DateInput,
  DescriptionInput,
  FormActionButton,
  HeaderAmountInput,
  NoteInput,
} from './FormFields.tsx'
import { AccountSelector, CategorySelector } from './Selectors.tsx'
import { AddCategoryForm } from './AddCategoryForm.tsx'
import { TransactionTypeSelector } from './TransactionTypeSelector.tsx'

type TransactionFormValues = {
  description: string
  amount: number
  occurredOn: string
  accountId: string
  categoryId: string
  sourceAccountId: string
  destinationAccountId: string
  note: string
}

type TransactionFormProps = {
  type: TransactionType
  mode?: 'create' | 'edit'
  initialTransaction?: Transaction
}

function getFormSchema(type: TransactionType) {
  return z.object({
    description: z.string().trim().min(1, 'Informe a descricao.'),
    amount: z.number().int().positive('Informe um valor maior que zero.'),
    occurredOn: z.string().min(1, 'Informe a data.'),
    accountId: z.string(),
    categoryId: z.string(),
    sourceAccountId: z.string(),
    destinationAccountId: z.string(),
    note: z.string(),
  }).superRefine((values, context) => {
    if (type === 'transfer') {
      if (!values.sourceAccountId) {
        context.addIssue({
          code: 'custom',
          path: ['sourceAccountId'],
          message: 'Selecione a conta de origem.',
        })
      }
      if (!values.destinationAccountId) {
        context.addIssue({
          code: 'custom',
          path: ['destinationAccountId'],
          message: 'Selecione a conta de destino.',
        })
      }
      if (
        values.sourceAccountId &&
        values.destinationAccountId &&
        values.sourceAccountId === values.destinationAccountId
      ) {
        context.addIssue({
          code: 'custom',
          path: ['destinationAccountId'],
          message: 'Origem e destino devem ser diferentes.',
        })
      }
      return
    }
    if (!values.accountId) {
      context.addIssue({
        code: 'custom',
        path: ['accountId'],
        message: 'Selecione uma conta.',
      })
    }
    if (!values.categoryId) {
      context.addIssue({
        code: 'custom',
        path: ['categoryId'],
        message: 'Selecione uma categoria.',
      })
    }
  })
}

function getDefaultValues(): TransactionFormValues {
  return {
    description: '',
    amount: 0,
    occurredOn: toDateInputValue(new Date()),
    accountId: '',
    categoryId: '',
    sourceAccountId: '',
    destinationAccountId: '',
    note: '',
  }
}

function getInitialValues(transaction: Transaction): TransactionFormValues {
  return {
    description: transaction.description,
    amount: transaction.amount,
    occurredOn: toDateInputValue(new Date(transaction.occurredAt)),
    accountId: transaction.accountId ?? '',
    categoryId: transaction.categoryId ?? '',
    sourceAccountId: transaction.sourceAccountId ?? '',
    destinationAccountId: transaction.destinationAccountId ?? '',
    note: transaction.note,
  }
}

function getSubmitLabel(type: TransactionType, mode: 'create' | 'edit') {
  if (mode === 'edit') {
    return 'Salvar alteracoes'
  }
  if (type === 'income') {
    return 'Salvar receita'
  }
  if (type === 'transfer') {
    return 'Salvar transferencia'
  }
  return 'Salvar despesa'
}

function toOccurredAt(value: string) {
  return formatLocalRFC3339(fromDateInputValue(value))
}

const typeStyles = {
  income: {
    header: 'bg-[#159c57]',
    createTitle: 'Nova receita',
    editTitle: 'Editar receita',
    amountLabel: 'Valor da receita',
    accentColor: '#159c57',
  },
  expense: {
    header: 'bg-[#d93658]',
    createTitle: 'Nova despesa',
    editTitle: 'Editar despesa',
    amountLabel: 'Valor da despesa',
    accentColor: '#d93658',
  },
  transfer: {
    header: 'bg-[#2478d4]',
    createTitle: 'Nova transferencia',
    editTitle: 'Editar transferencia',
    amountLabel: 'Valor da transferencia',
    accentColor: '#2478d4',
  },
} satisfies Record<
  TransactionType,
  {
    header: string
    createTitle: string
    editTitle: string
    amountLabel: string
    accentColor: string
  }
>

export function TransactionForm({
  type,
  mode = 'create',
  initialTransaction,
}: TransactionFormProps) {
  const navigate = useNavigate()
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const createTransactionMutation = useCreateTransaction(type)
  const updateTransactionMutation = useUpdateTransaction(
    type,
    initialTransaction?.id ?? '',
  )
  const resolver = useMemo(() => zodResolver(getFormSchema(type)), [type])
  const defaultValues = useMemo(
    () =>
      mode === 'edit' && initialTransaction
        ? getInitialValues(initialTransaction)
        : getDefaultValues(),
    [initialTransaction, mode],
  )
  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<TransactionFormValues>({
    defaultValues,
    resolver,
  })
  const categoryType = type === 'income' || type === 'expense' ? type : null
  const styles = typeStyles[type]
  const isPending =
    mode === 'edit'
      ? updateTransactionMutation.isPending
      : createTransactionMutation.isPending
  const hasMutationError =
    mode === 'edit'
      ? updateTransactionMutation.isError
      : createTransactionMutation.isError

  return (
    <>
      <form
        className="mx-auto min-h-svh w-full max-w-[520px] bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:max-w-[640px]"
        onSubmit={handleSubmit((values) => {
          const base = {
            description: values.description.trim(),
            amount: values.amount,
            occurredAt: toOccurredAt(values.occurredOn),
            note: values.note.trim(),
          }

          if (type === 'transfer') {
            const payload: CreateTransferTransactionPayload = {
              ...base,
              sourceAccountId: values.sourceAccountId,
              destinationAccountId: values.destinationAccountId,
            }
            if (mode === 'edit') {
              updateTransactionMutation.mutate(
                payload as UpdateTransactionPayloadByType[typeof type],
                {
                  onSuccess: () => navigate('/transactions'),
                },
              )
              return
            }
            createTransactionMutation.mutate(payload, {
              onSuccess: () => navigate('/transactions'),
            })
            return
          }

          const payload:
            | CreateIncomeTransactionPayload
            | CreateExpenseTransactionPayload = {
            ...base,
            accountId: values.accountId,
            categoryId: values.categoryId,
          }
          if (mode === 'edit') {
            updateTransactionMutation.mutate(
              payload as UpdateTransactionPayloadByType[typeof type],
              {
                onSuccess: () => navigate('/transactions'),
              },
            )
            return
          }
          createTransactionMutation.mutate(payload, {
            onSuccess: () => navigate('/transactions'),
          })
        })}
      >
        <div className={`${styles.header} w-full px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7 lg:px-10 lg:pb-14`}>
          <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center md:grid-cols-[minmax(0,1fr)_auto] md:gap-4">
            <button
              type="button"
              className="justify-self-start text-[14px] font-semibold text-white/88 transition-colors hover:text-white focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-white md:order-2"
              onClick={() => navigate('/transactions')}
            >
              Cancelar
            </button>
            <div className="min-w-0 md:order-1 md:justify-self-start">
              <TransactionTypeSelector type={type} isLocked={mode === 'edit'} />
            </div>
          </div>
          <h1 className="sr-only">
            {mode === 'edit' ? styles.editTitle : styles.createTitle}
          </h1>
          <div className="mt-10">
            <Controller
              control={control}
              name="amount"
              render={({ field }) => (
                <HeaderAmountInput
                  label={styles.amountLabel}
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
                value={field.value}
                accentColor={styles.accentColor}
                error={errors.occurredOn?.message}
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
          {type === 'transfer' ? (
            <>
              <Controller
                control={control}
                name="sourceAccountId"
                render={({ field }) => (
                  <AccountSelector
                    label="Origem"
                    value={field.value}
                    error={errors.sourceAccountId?.message}
                    onChange={field.onChange}
                  />
                )}
              />
              <Controller
                control={control}
                name="destinationAccountId"
                render={({ field }) => (
                  <AccountSelector
                    label="Destino"
                    value={field.value}
                    error={errors.destinationAccountId?.message}
                    onChange={field.onChange}
                  />
                )}
              />
            </>
          ) : (
            <>
              <Controller
                control={control}
                name="categoryId"
                render={({ field }) => (
                  <CategorySelector
                    type={type}
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
            </>
          )}
          <Controller
            control={control}
            name="note"
            render={({ field }) => <NoteInput value={field.value} onChange={field.onChange} />}
          />
          {hasMutationError ? (
            <p className="mx-4 mt-4 rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838] md:mx-5">
              Nao foi possivel salvar a transacao.
            </p>
          ) : null}
          <FormActionButton isPending={isPending}>
            {getSubmitLabel(type, mode)}
          </FormActionButton>
        </div>
      </form>

      {categoryType ? (
        <AddCategoryForm
          type={categoryType as CategoryTransactionType}
          isOpen={isAddingCategory}
          onClose={() => setIsAddingCategory(false)}
          onCreated={(categoryId) => setValue('categoryId', categoryId, { shouldValidate: true })}
        />
      ) : null}
    </>
  )
}
