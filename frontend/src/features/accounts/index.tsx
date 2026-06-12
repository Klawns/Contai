import { useCallback, useMemo } from 'react'
import { ArrowLeft, CircleDollarSign } from 'lucide-react'
import { motion, useReducedMotion } from 'motion/react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useConfirmDialog } from '../../components/confirm-dialog-context.ts'
import { MonthSelector, type SelectedMonth } from '../../components/MonthSelector.tsx'
import { ItemActionsMenu } from '../../components/ItemActionsMenu.tsx'
import { formatCurrency } from '../transactions/utils/money.ts'
import {
  formatMonthQuery,
  getCurrentSelectedMonth,
} from '../transactions/utils/date.ts'
import { useAccounts, useTotalBalance } from './hooks/useAccounts.ts'
import { useDeleteAccount } from './hooks/useSaveAccount.ts'
import type { Account } from './types/accounts.ts'
import { BankIcon } from './components/BankIcon.tsx'
import { AccountForm } from './components/AccountForm.tsx'

const monthQueryPattern = /^(\d{4})-(0[1-9]|1[0-2])$/

function parseMonthQuery(value: string | null): SelectedMonth | null {
  const match = value?.match(monthQueryPattern)

  if (!match) {
    return null
  }

  return {
    year: Number(match[1]),
    monthIndex: Number(match[2]) - 1,
  }
}

function getBalanceColor(valueInCents: number) {
  if (valueInCents > 0) {
    return 'text-[#18794e]'
  }

  if (valueInCents < 0) {
    return 'text-[#c83b3b]'
  }

  return 'text-[#2f263b]'
}

function AccountRow({ account }: { account: Account }) {
  const navigate = useNavigate()
  const { confirm } = useConfirmDialog()
  const shouldReduceMotion = useReducedMotion()
  const deleteAccountMutation = useDeleteAccount()
  const balanceColor = getBalanceColor(account.currentBalance)

  async function handleDeleteAccount() {
    const shouldDelete = await confirm({
      title: 'Deletar conta',
      description: `Deletar a conta "${account.name}"? Esta acao nao pode ser desfeita.`,
      confirmLabel: 'Deletar',
      cancelLabel: 'Cancelar',
      tone: 'danger',
    })

    if (!shouldDelete) {
      return
    }

    deleteAccountMutation.mutate(account.id)
  }

  return (
    <motion.li
      className="grid grid-cols-[34px_minmax(0,1fr)_minmax(86px,auto)_32px] items-center gap-3 px-1 py-2.5 sm:grid-cols-[36px_minmax(0,1fr)_minmax(112px,auto)_32px] sm:px-2"
      whileHover={shouldReduceMotion ? undefined : { x: 2 }}
      transition={{ duration: 0.16, ease: 'easeOut' }}
    >
      <BankIcon bankIconId={account.bankIconId} size={34} />
      <div className="min-w-0">
        <h3 className="truncate text-[14px] font-semibold leading-tight text-[#241a30]">
          {account.name}
        </h3>
        <p className="mt-1 truncate text-[12px] font-semibold leading-tight text-[#81788c]">
          Saldo atual
        </p>
      </div>
      <div className="min-w-0 text-right">
        <strong className={`block text-[13px] font-semibold leading-tight sm:text-[14px] ${balanceColor}`}>
          {formatCurrency(account.currentBalance)}
        </strong>
      </div>
      <ItemActionsMenu
        label={`Acoes de ${account.name}`}
        onEdit={() => navigate(`/accounts/edit?accountId=${encodeURIComponent(account.id)}`)}
        onDelete={handleDeleteAccount}
        isDeleteDisabled={deleteAccountMutation.isPending}
      />
    </motion.li>
  )
}

function AccountsList({
  accounts,
  isLoading,
  isError,
}: {
  accounts?: Account[]
  isLoading: boolean
  isError: boolean
}) {
  if (isLoading) {
    return <StatePanel>Carregando contas...</StatePanel>
  }

  if (isError) {
    return <StatePanel tone="danger">Nao foi possivel carregar as contas.</StatePanel>
  }

  if (!accounts?.length) {
    return (
      <StatePanel>
        <span>Ainda nao ha contas cadastradas.</span>
      </StatePanel>
    )
  }

  return (
    <section className="bg-white">
      <ul className="divide-y divide-[#f0ebf6]">
        {accounts.map((account) => (
          <AccountRow key={account.id} account={account} />
        ))}
      </ul>
    </section>
  )
}

function StatePanel({
  tone = 'default',
  children,
}: {
  tone?: 'default' | 'danger'
  children: React.ReactNode
}) {
  return (
    <div
      className={`grid justify-items-center rounded-2xl border bg-white px-5 py-5 text-center text-[14px] font-semibold ${
        tone === 'danger'
          ? 'border-[#f0caca] text-[#b93838]'
          : 'border-[#ece8f2] text-[#81788c]'
      }`}
    >
      {children}
    </div>
  )
}

function AccountStatePage({
  tone = 'default',
  children,
}: {
  tone?: 'default' | 'danger'
  children: React.ReactNode
}) {
  return (
    <main className="scrollbar-none h-full min-h-0 w-full overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left md:px-8">
      <StatePanel tone={tone}>{children}</StatePanel>
    </main>
  )
}

export function AccountListPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const selectedMonth = useMemo(
    () => parseMonthQuery(searchParams.get('month')) ?? getCurrentSelectedMonth(),
    [searchParams],
  )
  const accountsQuery = useAccounts()
  const totalBalanceQuery = useTotalBalance()
  const updateSelectedMonth = useCallback((nextMonth: SelectedMonth) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set('month', formatMonthQuery(nextMonth))
    setSearchParams(nextParams)
  }, [searchParams, setSearchParams])

  return (
    <main className="h-full min-h-0 w-full min-w-0 overflow-hidden bg-[#6818e8] text-left" aria-label="Contas">
      <section className="mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-hidden bg-[#6818e8] shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none">
        <header className="flex-none bg-[#6818e8] px-5 pb-5 pt-[calc(18px+env(safe-area-inset-top))] text-white md:px-7 md:pt-6">
          <div className="mx-auto grid w-full grid-cols-[44px_minmax(0,1fr)_44px] items-center">
            <button
              type="button"
              className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              aria-label="Voltar"
              onClick={() => navigate('/')}
            >
              <ArrowLeft className="h-5 w-5" aria-hidden="true" />
            </button>
            <div className="min-w-0 px-2 text-center">
              <h1 className="truncate text-[17px] font-semibold leading-tight md:text-[24px]">
                Contas
              </h1>
              <div className="mt-1 flex justify-center">
                <MonthSelector
                  selectedMonth={selectedMonth}
                  onSelectMonth={updateSelectedMonth}
                  tone="inverse"
                />
              </div>
            </div>
            <div aria-hidden="true" />
          </div>
        </header>

        <div className="scrollbar-none flex min-h-0 flex-1 flex-col overflow-y-auto overflow-x-hidden rounded-t-[26px] bg-white px-5 pb-[var(--app-mobile-content-bottom)] pt-4 md:px-7 md:pb-10">
          <div className="flex w-full min-w-0 flex-1 flex-col gap-2.5">
            <section className="border-b border-[#f0ebf6] pb-4">
              <div className="min-w-0">
                <CircleDollarSign className="h-5 w-5 text-[#1f9d63]" aria-hidden="true" />
                <span className="mt-2 block text-[12px] font-semibold leading-tight text-[#81788c]">
                  Saldo atual
                </span>
                <strong className="mt-1 block truncate text-[18px] font-semibold leading-tight text-[#18794e] sm:text-[20px]">
                  {totalBalanceQuery.isLoading
                    ? '...'
                    : formatCurrency(totalBalanceQuery.data?.totalBalance ?? 0)}
                </strong>
              </div>
            </section>

            <AccountsList
              accounts={accountsQuery.data}
              isLoading={accountsQuery.isLoading}
              isError={accountsQuery.isError}
            />

            {!accountsQuery.isLoading && !accountsQuery.isError ? (
              <div className="flex justify-center pt-1">
                <button
                  type="button"
                  className="min-h-9 cursor-pointer bg-transparent px-4 text-[14px] font-semibold text-[#6a22e5] transition-colors hover:text-[#5114bd] focus-visible:rounded-full focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#8f57ff]"
                  onClick={() => navigate('/accounts/new')}
                >
                  Cadastrar conta
                </button>
              </div>
            ) : null}
          </div>
        </div>
      </section>
    </main>
  )
}

export function AccountCreatePage() {
  return <AccountForm mode="create" />
}

export function AccountEditPage() {
  const [searchParams] = useSearchParams()
  const accountId = searchParams.get('accountId')
  const accountsQuery = useAccounts()
  const account = accountsQuery.data?.find((item) => item.id === accountId)

  if (!accountId) {
    return <AccountStatePage tone="danger">Conta nao informada.</AccountStatePage>
  }

  if (accountsQuery.isLoading) {
    return <AccountStatePage>Carregando conta...</AccountStatePage>
  }

  if (accountsQuery.isError || !account) {
    return <AccountStatePage tone="danger">Nao foi possivel encontrar a conta.</AccountStatePage>
  }

  return <AccountForm mode="edit" account={account} />
}
