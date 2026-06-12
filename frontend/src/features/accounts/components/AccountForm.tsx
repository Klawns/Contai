import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import type { ReactNode } from 'react'
import { Check, ChevronRight, Info, Landmark, PencilLine } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { HeaderAmountInput } from '../../transactions/components/FormFields.tsx'
import { SelectionSheet } from '../../transactions/components/SelectionSheet.tsx'
import { formatCurrency } from '../../transactions/utils/money.ts'
import {
  createAccountPayloadSchema,
} from '../schemas/accounts.ts'
import type { Account, AccountType, CreateAccountPayload, UpdateAccountPayload } from '../types/accounts.ts'
import { accountTypeLabels } from '../utils/formatters.ts'
import { BankIconSelector } from './BankIconSelector.tsx'
import { useCreateAccount, useUpdateAccount } from '../hooks/useSaveAccount.ts'

type AccountFormValues = {
  name: string
  type: AccountType
  initialBalance: number
  bankIconId: string
  includeInDashboardTotal: boolean
}

type AccountFormProps = {
  mode: 'create' | 'edit'
  account?: Account
}

const accountTypeOptions = Object.entries(accountTypeLabels).map(([value, label]) => ({
  value: value as AccountType,
  label,
}))

function getDefaultValues(account?: Account): AccountFormValues {
  return {
    name: account?.name ?? '',
    type: account?.type ?? 'checking',
    initialBalance: account?.initialBalance ?? 0,
    bankIconId: account?.bankIconId ?? '',
    includeInDashboardTotal: account?.includeInDashboardTotal ?? true,
  }
}

export function AccountForm({ mode, account }: AccountFormProps) {
  const navigate = useNavigate()
  const isEditing = mode === 'edit'
  const createAccountMutation = useCreateAccount()
  const updateAccountMutation = useUpdateAccount(account?.id ?? '')
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<AccountFormValues>({
    defaultValues: getDefaultValues(account),
    resolver: zodResolver(createAccountPayloadSchema),
  })
  const isPending = createAccountMutation.isPending || updateAccountMutation.isPending
  const isError = createAccountMutation.isError || updateAccountMutation.isError

  return (
    <form
      className="scrollbar-none mx-auto h-full min-h-0 w-full max-w-[520px] overflow-y-auto overflow-x-hidden bg-white text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:flex md:max-w-none md:flex-col md:shadow-none"
      onSubmit={handleSubmit((values) => {
        if (isEditing && account) {
          const payload: UpdateAccountPayload = {
            name: values.name.trim(),
            type: values.type,
            bankIconId: values.bankIconId,
            includeInDashboardTotal: values.includeInDashboardTotal,
          }
          updateAccountMutation.mutate(payload, {
            onSuccess: () => navigate('/accounts'),
          })
          return
        }

        const payload: CreateAccountPayload = {
          name: values.name.trim(),
          type: values.type,
          initialBalance: values.initialBalance,
          bankIconId: values.bankIconId,
          includeInDashboardTotal: values.includeInDashboardTotal,
        }
        createAccountMutation.mutate(payload, {
          onSuccess: () => navigate('/accounts'),
        })
      })}
    >
      <div className="w-full bg-[#6818e8] px-5 pb-14 pt-[calc(22px+env(safe-area-inset-top))] text-white md:px-8 md:pb-12 md:pt-7 lg:px-10 lg:pb-14">
        <div className="grid grid-cols-[80px_minmax(0,1fr)_80px] items-center md:grid-cols-[minmax(0,1fr)_auto] md:gap-4">
          <button
            type="button"
            className="justify-self-start text-[14px] font-semibold text-white/88 transition-colors hover:text-white focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-white md:order-2"
            onClick={() => navigate('/accounts')}
          >
            Cancelar
          </button>
          <h1 className="truncate text-center text-[17px] font-semibold leading-tight md:order-1 md:text-left md:text-[24px]">
            {isEditing ? 'Editar Conta' : 'Nova Conta'}
          </h1>
        </div>

        <div className="mt-10">
          {isEditing ? (
            <div className="grid gap-2 border-b border-white/24 pb-3">
              <span className="text-left text-[14px] font-semibold text-white/78">
                Saldo atual
              </span>
              <strong className="block text-[48px] font-bold leading-none text-white sm:text-[64px]">
                {formatCurrency(account?.currentBalance ?? 0)}
              </strong>
            </div>
          ) : (
            <Controller
              control={control}
              name="initialBalance"
              render={({ field }) => (
                <HeaderAmountInput
                  label="Saldo inicial"
                  value={field.value}
                  error={errors.initialBalance?.message}
                  onChange={field.onChange}
                />
              )}
            />
          )}
        </div>
      </div>

      <div className="-mt-6 rounded-t-[28px] border border-[#ece8f2] bg-white md:flex md:w-full md:flex-1 md:flex-col md:px-8 md:py-7 lg:px-10">
        <Controller
          control={control}
          name="bankIconId"
          render={({ field }) => (
            <div className="md:col-span-2">
              <BankIconSelector
                value={field.value}
                error={errors.bankIconId?.message}
                onChange={field.onChange}
              />
            </div>
          )}
        />

        <Controller
          control={control}
          name="name"
          render={({ field }) => (
            <AccountFieldRow
              icon={<PencilLine className="h-5 w-5" aria-hidden="true" />}
              error={errors.name?.message}
            >
              <input
                className="h-12 w-full bg-transparent text-left text-[16px] font-normal text-[#2c2237] outline-none placeholder:text-[#aaa2b4]"
                value={field.value}
                placeholder="Nome da conta"
                onChange={field.onChange}
              />
            </AccountFieldRow>
          )}
        />

        <Controller
          control={control}
          name="type"
          render={({ field }) => (
            <AccountTypeSelector
              value={field.value}
              error={errors.type?.message}
              onChange={field.onChange}
            />
          )}
        />

        <Controller
          control={control}
          name="includeInDashboardTotal"
          render={({ field }) => (
            <AccountSwitchRow
              checked={field.value}
              onChange={field.onChange}
            />
          )}
        />

        {isError ? (
          <p className="mx-4 mt-3 rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838] md:col-span-2 md:mx-0 md:mt-0">
            Nao foi possivel salvar a conta.
          </p>
        ) : null}

        <div className="sticky bottom-[var(--app-mobile-sticky-bottom)] mt-4 border-t border-[#eee8f3] bg-white/96 px-4 py-4 backdrop-blur md:static md:border-t-0 md:bg-transparent md:px-0 md:pt-6 md:backdrop-blur-none">
          <button
            type="submit"
            disabled={isPending}
            className="mx-auto flex h-12 w-full max-w-[300px] cursor-pointer items-center justify-center rounded-full bg-[#6818e8] px-6 text-[15px] font-semibold text-white shadow-[0_6px_14px_rgba(104,24,232,0.12)] transition-colors hover:bg-[#5712c9] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] disabled:cursor-not-allowed disabled:opacity-65"
          >
            {isPending ? 'Salvando...' : isEditing ? 'Salvar alteracoes' : 'Cadastrar conta'}
          </button>
        </div>
      </div>
    </form>
  )
}

type AccountTypeSelectorProps = {
  value: AccountType
  error?: string
  onChange: (value: AccountType) => void
}

function AccountTypeSelector({ value, error, onChange }: AccountTypeSelectorProps) {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <AccountFieldRow
        icon={<Landmark className="h-5 w-5" aria-hidden="true" />}
        error={error}
        trailing={<ChevronRight className="h-5 w-5" aria-hidden="true" />}
      >
        <button
          type="button"
          className="h-12 w-full cursor-pointer truncate bg-transparent text-left text-[16px] font-normal text-[#2c2237] outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
          aria-label="Tipo de conta"
          onClick={() => setIsOpen(true)}
        >
          {accountTypeLabels[value]}
        </button>
      </AccountFieldRow>

      <SelectionSheet title="Tipo de conta" isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="grid gap-2">
          {accountTypeOptions.map((option) => {
            const isSelected = option.value === value

            return (
              <button
                key={option.value}
                type="button"
                className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
                  isSelected
                    ? 'border-[#7b2cff] bg-[#f7f2ff]'
                    : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
                }`}
                aria-pressed={isSelected}
                onClick={() => {
                  onChange(option.value)
                  setIsOpen(false)
                }}
              >
                <span className="grid h-10 w-10 flex-none place-items-center rounded-full bg-[#f0ebf8] text-[#6a22e5]">
                  <Landmark className="h-5 w-5" aria-hidden="true" />
                </span>
                <span className="min-w-0 flex-1 truncate text-[14px] font-semibold text-[#2c2237]">
                  {option.label}
                </span>
                {isSelected ? (
                  <Check className="h-5 w-5 flex-none text-[#6818e8]" aria-hidden="true" />
                ) : null}
              </button>
            )
          })}
        </div>
      </SelectionSheet>
    </>
  )
}

type AccountSwitchRowProps = {
  checked: boolean
  onChange: (checked: boolean) => void
}

function AccountSwitchRow({ checked, onChange }: AccountSwitchRowProps) {
  return (
    <label className="relative block min-h-[68px] cursor-pointer border-b border-[#f0ebf5] px-4 py-2.5 md:px-5">
      <span className="absolute left-4 top-1/2 z-10 grid h-10 w-10 -translate-y-1/2 place-items-center text-[#6f647b] md:left-5">
        <Info className="h-5 w-5" aria-hidden="true" />
      </span>
      <span className="pointer-events-none absolute left-[68px] right-[84px] top-0 flex h-full items-center overflow-hidden md:left-[72px] md:right-[88px]">
        <span className="truncate text-left text-[16px] font-normal leading-snug text-[#2c2237]">
          Incluir na soma da tela inicial
        </span>
      </span>
      <span className="absolute right-4 top-1/2 z-10 h-8 w-[52px] -translate-y-1/2 md:right-5">
        <input
          type="checkbox"
          className="peer sr-only"
          checked={checked}
          aria-label="Incluir na soma da tela inicial"
          onChange={(event) => onChange(event.target.checked)}
        />
        <span className="absolute inset-0 rounded-full bg-[#d8d3df] transition-colors peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-[#7b2cff] peer-checked:bg-[#6818e8]" aria-hidden="true" />
        <span className="absolute left-1 top-1 h-6 w-6 rounded-full bg-white shadow-[0_2px_6px_rgba(43,35,54,0.22)] transition-transform peer-checked:translate-x-5" aria-hidden="true" />
      </span>
    </label>
  )
}

type AccountFieldRowProps = {
  icon: ReactNode
  error?: string
  className?: string
  trailing?: ReactNode
  children: ReactNode
}

function AccountFieldRow({
  icon,
  error,
  className = '',
  trailing,
  children,
}: AccountFieldRowProps) {
  return (
    <div className={`border-b border-[#f0ebf5] px-4 py-2.5 md:px-5 ${className}`}>
      <div className="relative min-h-12">
        <span className="absolute left-0 top-1/2 z-10 grid h-10 w-10 -translate-y-1/2 place-items-center text-[#6f647b]">
          {icon}
        </span>
        <div className="absolute left-[52px] right-[52px] top-0 flex h-full items-center overflow-hidden">
          <div className="w-full min-w-0 text-left">
            {children}
          </div>
        </div>
        <span className="absolute right-0 top-1/2 z-10 grid h-10 w-[52px] -translate-y-1/2 place-items-center text-[#9a91a5]">
          {trailing}
        </span>
      </div>
      {error ? (
        <p className="mt-1 text-center text-[12px] font-semibold text-[#c72f4d]">
          {error}
        </p>
      ) : null}
    </div>
  )
}
