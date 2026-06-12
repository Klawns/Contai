import { useCallback, useEffect, useId, useMemo, useRef, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  ArrowLeft,
  ArrowLeftRight,
  CalendarDays,
  ChevronLeft,
  ChevronRight,
  CircleDollarSign,
  Ellipsis,
  Landmark,
  Tag,
  TrendingDown,
  TrendingUp,
} from 'lucide-react'
import { useConfirmDialog } from '../../components/confirm-dialog-context.ts'
import {
  MonthSelector,
  type SelectedMonth,
} from '../../components/MonthSelector.tsx'
import { ItemActionsMenu } from '../../components/ItemActionsMenu.tsx'
import { useActiveAccounts } from './hooks/useActiveAccounts.ts'
import { useActiveCategories } from './hooks/useActiveCategories.ts'
import { useTransactions } from './hooks/useTransactions.ts'
import { useDeleteTransaction } from './hooks/useDeleteTransaction.ts'
import type { Category, Transaction, TransactionType } from './types/transactions.ts'
import {
  formatMonthQuery,
  getCurrentSelectedMonth,
  getMonthPeriod,
} from './utils/date.ts'
import { formatCurrency } from './utils/money.ts'
import { TransactionForm } from './components/TransactionForm.tsx'
import { TransactionStateMessage } from './components/TransactionStateMessage.tsx'
import { TransactionsPageLayout } from './components/TransactionsPageLayout.tsx'
import { useTotalBalance } from '../accounts/hooks/useAccounts.ts'

const createOptions = [
  {
    type: 'income',
    label: 'Receita',
    path: '/transactions/income/new',
    icon: TrendingUp,
    tone: 'bg-[#e8f8ef] text-[#147a46]',
  },
  {
    type: 'expense',
    label: 'Despesa',
    path: '/transactions/expense/new',
    icon: TrendingDown,
    tone: 'bg-[#fff0f2] text-[#c72f4d]',
  },
  {
    type: 'transfer',
    label: 'Transferencia',
    path: '/transactions/transfer/new',
    icon: ArrowLeftRight,
    tone: 'bg-[#eef6ff] text-[#216fb8]',
  },
] as const satisfies readonly {
  type: TransactionType
  label: string
  path: string
  icon: typeof TrendingUp
  tone: string
}[]

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

function getSignedAmount(transaction: Transaction) {
  if (transaction.settlementStatus !== 'settled') {
    return 0
  }
  if (transaction.type === 'expense') {
    return -transaction.amount
  }
  if (transaction.type === 'income') {
    return transaction.amount
  }
  return 0
}

function affectsBalance(transaction: Transaction) {
  if (transaction.type === 'transfer') {
    return true
  }
  return transaction.settlementStatus === 'settled' && Boolean(transaction.accountId)
}

function getSettlementLabel(transaction: Transaction) {
  if (transaction.type === 'transfer') {
    return null
  }
  if (transaction.type === 'income') {
    return transaction.settlementStatus === 'settled' ? 'Recebido' : 'Nao recebido'
  }
  return transaction.settlementStatus === 'settled' ? 'Pago' : 'Nao pago'
}

function getTransactionIcon(transaction: Transaction) {
  if (transaction.type === 'income') {
    return TrendingUp
  }
  if (transaction.type === 'expense') {
    return TrendingDown
  }
  return ArrowLeftRight
}

function getTransactionTone(transaction: Transaction) {
  if (transaction.type === 'income') {
    return 'bg-[#e8f8ef] text-[#147a46]'
  }
  if (transaction.type === 'expense') {
    return 'bg-[#fff0f2] text-[#c72f4d]'
  }
  return 'bg-[#eef6ff] text-[#216fb8]'
}

function getTransactionOriginLabel(transaction: Transaction) {
  if (transaction.originType === 'payable') {
    return 'Gerada por conta a pagar'
  }
  if (transaction.originType === 'receivable') {
    return 'Gerada por conta a receber'
  }
  if (transaction.originType === 'credit_card_invoice') {
    return 'Gerada por fatura de cartao'
  }
  return null
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'short',
  }).format(new Date(value))
}

function toMap<TItem extends { id: string }>(items: TItem[] | undefined) {
  return new Map((items ?? []).map((item) => [item.id, item]))
}

function getCategoryMap(incomeCategories?: Category[], expenseCategories?: Category[]) {
  return new Map(
    [...(incomeCategories ?? []), ...(expenseCategories ?? [])].map((category) => [
      category.id,
      category,
    ]),
  )
}

function CompactTransactionSummary({
  monthlyBalance,
  totalBalance,
  isTotalBalanceLoading,
}: {
  monthlyBalance: number
  totalBalance: number
  isTotalBalanceLoading: boolean
}) {
  return (
    <section className="grid grid-cols-2 border-b border-[#f0ebf6] pb-4">
      <div className="min-w-0 pr-4">
        <CircleDollarSign className="h-5 w-5 text-[#1f9d63]" aria-hidden="true" />
        <span className="mt-2 block text-[12px] font-semibold leading-tight text-[#81788c]">
          Saldo atual
        </span>
        <strong className="mt-1 block truncate text-[17px] font-semibold leading-tight text-[#18794e] sm:text-[19px]">
          {isTotalBalanceLoading ? '...' : formatCurrency(totalBalance)}
        </strong>
      </div>
      <div className="min-w-0 border-l border-[#f0ebf6] pl-4">
        <TrendingUp className="h-5 w-5 text-[#1f9d63]" aria-hidden="true" />
        <span className="mt-2 block text-[12px] font-semibold leading-tight text-[#81788c]">
          Balanco mensal
        </span>
        <strong className="mt-1 block truncate text-[17px] font-semibold leading-tight text-[#18794e] sm:text-[19px]">
          {formatCurrency(monthlyBalance)}
        </strong>
      </div>
    </section>
  )
}

function CreateTransactionMenu() {
  const navigate = useNavigate()
  const [isOpen, setIsOpen] = useState(false)
  const menuId = useId()
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!isOpen) {
      return undefined
    }

    function handlePointerDown(event: PointerEvent) {
      if (!menuRef.current?.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('pointerdown', handlePointerDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen])

  return (
    <div ref={menuRef} className="relative z-20">
      <button
        type="button"
        className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
        aria-label="Criar transacao"
        aria-haspopup="menu"
        aria-expanded={isOpen}
        aria-controls={menuId}
        onClick={() => setIsOpen((current) => !current)}
      >
        <Ellipsis className="h-5 w-5" aria-hidden="true" />
      </button>

      {isOpen ? (
        <div
          id={menuId}
          role="menu"
          className="absolute right-0 top-12 w-[210px] overflow-hidden rounded-2xl bg-white py-1.5 text-[#2c2237] shadow-[0_18px_45px_rgba(35,24,52,0.24)] ring-1 ring-black/6"
        >
          {createOptions.map((option) => {
            const Icon = option.icon

            return (
              <button
                key={option.type}
                type="button"
                role="menuitem"
                className="flex h-12 w-full cursor-pointer items-center gap-3 px-3.5 text-left text-[14px] font-semibold text-[#4f435c] transition-colors hover:bg-[#f8f5fb]"
                onClick={() => {
                  setIsOpen(false)
                  navigate(option.path)
                }}
              >
                <span
                  className={`grid h-8 w-8 flex-none place-items-center rounded-full ${option.tone}`}
                >
                  <Icon className="h-4 w-4" aria-hidden="true" />
                </span>
                <span>{option.label}</span>
              </button>
            )
          })}
        </div>
      ) : null}
    </div>
  )
}

function TransactionActionsMenu({ transaction }: { transaction: Transaction }) {
  const navigate = useNavigate()
  const { confirm } = useConfirmDialog()
  const deleteTransactionMutation = useDeleteTransaction()
  const isManaged = transaction.originType !== 'manual'

  if (isManaged) {
    return (
      <span
        className="grid h-8 w-8 place-items-center rounded-full text-[#b2a9bd]"
        title="Transacao gerenciada"
        aria-label="Transacao gerenciada"
      >
        <CalendarDays className="h-4 w-4" aria-hidden="true" />
      </span>
    )
  }

  function handleEditTransaction() {
    navigate(`/transactions/edit?transactionId=${encodeURIComponent(transaction.id)}`)
  }

  async function handleDeleteTransaction() {
    const shouldDelete = await confirm({
      title: 'Deletar transacao',
      description: `Deletar a transacao "${transaction.description}"? Esta acao nao pode ser desfeita.`,
      confirmLabel: 'Deletar',
      cancelLabel: 'Cancelar',
      tone: 'danger',
    })

    if (!shouldDelete) {
      return
    }

    deleteTransactionMutation.mutate(transaction.id)
  }

  return (
    <ItemActionsMenu
      label={`Acoes de ${transaction.description}`}
      onEdit={handleEditTransaction}
      onDelete={handleDeleteTransaction}
      isDeleteDisabled={deleteTransactionMutation.isPending}
    />
  )
}

type TransactionListProps = {
  transactions: Transaction[]
  accountNames: Map<string, { name: string }>
  categoryNames: Map<string, { name: string }>
}

function TransactionList({ transactions, accountNames, categoryNames }: TransactionListProps) {
  if (transactions.length === 0) {
    return (
      <div className="grid flex-1 place-items-center px-5 py-12 text-center">
        <div className="grid justify-items-center">
          <div className="relative h-24 w-24" aria-hidden="true">
            <div className="absolute inset-x-5 bottom-4 h-12 rounded-2xl bg-[#f3ecff]" />
            <div className="absolute left-4 top-5 h-14 w-14 rounded-2xl border border-[#dbcdf3] bg-white shadow-[0_10px_24px_rgba(43,35,54,0.08)]" />
            <CalendarDays className="absolute left-8 top-9 h-8 w-8 text-[#6a22e5]" />
            <span className="absolute right-5 top-7 h-3 w-3 rounded-full bg-[#1f9d63]" />
            <span className="absolute bottom-7 right-7 h-2 w-7 rounded-full bg-[#d9c9f2]" />
          </div>
          <h2 className="mt-2 max-w-[280px] text-[17px] font-semibold leading-snug text-[#2c2237]">
            Ops, você não possui transações registradas
        </h2>
          <p className="mt-2 max-w-[300px] text-[13px] font-medium leading-relaxed text-[#81788c]">
            Adicione sua primeira receita, despesa ou transferencia pelo botao +.
          </p>
        </div>
      </div>
    )
  }

  return (
    <section className="bg-white">
      <ul className="divide-y divide-[#f0ebf6]">
        {transactions.map((transaction) => {
          const Icon = getTransactionIcon(transaction)
          const category = transaction.categoryId
            ? categoryNames.get(transaction.categoryId)
            : null
          const account = transaction.accountId
            ? accountNames.get(transaction.accountId)
            : null
          const source = transaction.sourceAccountId
            ? accountNames.get(transaction.sourceAccountId)
            : null
          const destination = transaction.destinationAccountId
            ? accountNames.get(transaction.destinationAccountId)
            : null
          const details =
            transaction.type === 'transfer'
              ? [source?.name, destination?.name].filter(Boolean).join(' -> ')
              : [account?.name, category?.name].filter(Boolean).join(' / ')
          const signedAmount = getSignedAmount(transaction)
          const originLabel = getTransactionOriginLabel(transaction)
          const settlementLabel = getSettlementLabel(transaction)
          const balanceLabel = !affectsBalance(transaction)
            ? transaction.accountId
              ? 'Pendente, nao movimentou saldo'
              : 'Sem conta, nao movimentou saldo'
            : null

          return (
            <li
              key={transaction.id}
              className="grid grid-cols-[40px_minmax(0,1fr)_minmax(82px,auto)_32px] items-center gap-3 px-1 py-3 sm:grid-cols-[40px_minmax(0,1fr)_minmax(112px,auto)_32px] sm:px-2"
            >
              <span
                className={`grid h-10 w-10 flex-none place-items-center rounded-full ${getTransactionTone(transaction)}`}
              >
                <Icon className="h-4.5 w-4.5" aria-hidden="true" />
              </span>
              <div className="min-w-0">
                <h3 className="truncate text-[14px] font-semibold leading-tight text-[#2c2237]">
                  {transaction.description}
                </h3>
                <p className="mt-1 flex min-w-0 items-center gap-1 text-[12px] font-semibold leading-tight text-[#81788c]">
                  {transaction.type === 'transfer' ? (
                    <Landmark className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
                  ) : (
                    <Tag className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
                  )}
                  <span className="truncate">{details || formatDate(transaction.occurredAt)}</span>
                </p>
                {originLabel ? (
                  <p className="mt-1 flex min-w-0 items-center gap-1 text-[11px] font-semibold leading-tight text-[#958c9f]">
                    <CalendarDays className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
                    <span className="truncate">{originLabel}</span>
                  </p>
                ) : null}
                {settlementLabel || balanceLabel ? (
                  <p className="mt-1 flex min-w-0 items-center gap-1 text-[11px] font-semibold leading-tight text-[#958c9f]">
                    <span className="truncate">
                      {[settlementLabel, balanceLabel].filter(Boolean).join(' - ')}
                    </span>
                  </p>
                ) : null}
              </div>
              <div className="min-w-0 text-right">
                <strong
                  className={`block text-[13px] font-semibold leading-tight sm:text-[14px] ${
                    signedAmount < 0
                      ? 'text-[#c72f4d]'
                      : signedAmount > 0
                        ? 'text-[#147a46]'
                        : 'text-[#216fb8]'
                  }`}
                >
                  {transaction.type === 'transfer'
                    ? formatCurrency(transaction.amount)
                    : formatCurrency(signedAmount)}
                </strong>
                <span className="mt-1 block text-[12px] font-semibold leading-tight text-[#958c9f]">
                  {formatDate(transaction.occurredAt)}
                </span>
              </div>
              <TransactionActionsMenu transaction={transaction} />
            </li>
          )
        })}
      </ul>
    </section>
  )
}

export function TransactionListPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const selectedMonth = useMemo(
    () => parseMonthQuery(searchParams.get('month')) ?? getCurrentSelectedMonth(),
    [searchParams],
  )
  const period = useMemo(() => getMonthPeriod(selectedMonth), [selectedMonth])
  const transactionsQuery = useTransactions(period)
  const totalBalanceQuery = useTotalBalance()
  const accountsQuery = useActiveAccounts()
  const incomeCategoriesQuery = useActiveCategories('income')
  const expenseCategoriesQuery = useActiveCategories('expense')
  const accountNames = useMemo(() => toMap(accountsQuery.data), [accountsQuery.data])
  const categoryNames = useMemo(
    () => getCategoryMap(incomeCategoriesQuery.data, expenseCategoriesQuery.data),
    [expenseCategoriesQuery.data, incomeCategoriesQuery.data],
  )
  const updateSelectedMonth = useCallback((nextMonth: SelectedMonth) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set('month', formatMonthQuery(nextMonth))
    setSearchParams(nextParams)
  }, [searchParams, setSearchParams])
  const changeSelectedMonth = useCallback(
    (offset: -1 | 1) => {
      const nextDate = new Date(selectedMonth.year, selectedMonth.monthIndex + offset, 1)

      updateSelectedMonth({
        year: nextDate.getFullYear(),
        monthIndex: nextDate.getMonth(),
      })
    },
    [selectedMonth, updateSelectedMonth],
  )
  const monthlyBalance = useMemo(
    () =>
      (transactionsQuery.data ?? []).reduce(
        (total, transaction) => total + getSignedAmount(transaction),
        0,
      ),
    [transactionsQuery.data],
  )

  return (
    <TransactionsPageLayout animationKey={formatMonthQuery(selectedMonth)}>
      <section className="mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-hidden bg-[#6818e8] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none">
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
                Transacoes
              </h1>
              <div className="mt-1 flex min-w-0 items-center justify-center gap-1">
                <button
                  type="button"
                  className="grid h-8 w-8 flex-none cursor-pointer place-items-center rounded-full text-white/88 transition-colors hover:bg-white/12 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                  aria-label="Mes anterior"
                  onClick={() => changeSelectedMonth(-1)}
                >
                  <ChevronLeft className="h-4.5 w-4.5" aria-hidden="true" />
                </button>
                <MonthSelector
                  selectedMonth={selectedMonth}
                  onSelectMonth={updateSelectedMonth}
                  tone="inverse"
                />
                <button
                  type="button"
                  className="grid h-8 w-8 flex-none cursor-pointer place-items-center rounded-full text-white/88 transition-colors hover:bg-white/12 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                  aria-label="Proximo mes"
                  onClick={() => changeSelectedMonth(1)}
                >
                  <ChevronRight className="h-4.5 w-4.5" aria-hidden="true" />
                </button>
              </div>
            </div>
            <CreateTransactionMenu />
          </div>
        </header>

        <div
          key={formatMonthQuery(selectedMonth)}
          className="scrollbar-none flex min-h-0 flex-1 flex-col overflow-y-auto overflow-x-hidden rounded-t-[26px] bg-white px-5 pb-[var(--app-mobile-content-bottom)] pt-4 md:px-7 md:pb-10"
        >
          <div className="flex w-full min-w-0 flex-1 flex-col gap-2.5">
            {transactionsQuery.isLoading ? (
              <TransactionStateMessage>Carregando transacoes...</TransactionStateMessage>
            ) : null}
            {transactionsQuery.isError ? (
              <TransactionStateMessage tone="danger">
                Nao foi possivel carregar as transacoes.
              </TransactionStateMessage>
            ) : null}
            {transactionsQuery.data ? (
              <>
                <CompactTransactionSummary
                  monthlyBalance={monthlyBalance}
                  totalBalance={totalBalanceQuery.data?.totalBalance ?? 0}
                  isTotalBalanceLoading={totalBalanceQuery.isLoading}
                />
                <TransactionList
                  transactions={transactionsQuery.data}
                  accountNames={accountNames}
                  categoryNames={categoryNames}
                />
              </>
            ) : null}
          </div>
        </div>
      </section>
    </TransactionsPageLayout>
  )
}

export function TransactionCreatePage({ type }: { type: TransactionType }) {
  return (
    <TransactionsPageLayout variant="create" tone={type} animationKey={type}>
      <TransactionForm type={type} />
    </TransactionsPageLayout>
  )
}

export function TransactionEditPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const transactionId = searchParams.get('transactionId') ?? ''
  const transactionsQuery = useTransactions({})
  const transaction = useMemo(
    () =>
      transactionId
        ? transactionsQuery.data?.find((item) => item.id === transactionId)
        : undefined,
    [transactionId, transactionsQuery.data],
  )

  if (!transactionId) {
    return (
      <TransactionsPageLayout variant="create" tone="expense" animationKey="edit-missing">
        <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:px-8 md:shadow-none">
          <button
            type="button"
            className="self-start text-[14px] font-semibold text-[#6818e8] transition-colors hover:text-[#4d12b0] focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#6818e8]"
            onClick={() => navigate('/transactions')}
          >
            Voltar
          </button>
          <TransactionStateMessage tone="danger">
            Informe uma transacao para editar.
          </TransactionStateMessage>
        </section>
      </TransactionsPageLayout>
    )
  }

  if (transactionsQuery.isLoading) {
    return (
      <TransactionsPageLayout variant="create" tone="expense" animationKey="edit-loading">
        <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:px-8 md:shadow-none">
          <TransactionStateMessage>Carregando transacao...</TransactionStateMessage>
        </section>
      </TransactionsPageLayout>
    )
  }

  if (transactionsQuery.isError) {
    return (
      <TransactionsPageLayout variant="create" tone="expense" animationKey="edit-error">
        <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:px-8 md:shadow-none">
          <button
            type="button"
            className="self-start text-[14px] font-semibold text-[#6818e8] transition-colors hover:text-[#4d12b0] focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#6818e8]"
            onClick={() => navigate('/transactions')}
          >
            Voltar
          </button>
          <TransactionStateMessage tone="danger">
            Nao foi possivel carregar a transacao.
          </TransactionStateMessage>
        </section>
      </TransactionsPageLayout>
    )
  }

  if (!transaction) {
    return (
      <TransactionsPageLayout variant="create" tone="expense" animationKey="edit-not-found">
        <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:px-8 md:shadow-none">
          <button
            type="button"
            className="self-start text-[14px] font-semibold text-[#6818e8] transition-colors hover:text-[#4d12b0] focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#6818e8]"
            onClick={() => navigate('/transactions')}
          >
            Voltar
          </button>
          <TransactionStateMessage tone="danger">
            Transacao nao encontrada.
          </TransactionStateMessage>
        </section>
      </TransactionsPageLayout>
    )
  }

  return (
    <TransactionsPageLayout
      variant="create"
      tone={transaction.type}
      animationKey={`edit-${transaction.id}`}
    >
      <TransactionForm
        type={transaction.type}
        mode="edit"
        initialTransaction={transaction}
      />
    </TransactionsPageLayout>
  )
}
