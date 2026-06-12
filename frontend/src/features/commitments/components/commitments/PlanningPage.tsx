import { useCallback, useMemo, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import type { SelectedMonth } from '../../../../components/MonthSelector.tsx'
import { mapById } from '../../../../lib/collections/mapById.ts'
import { TransactionStateMessage } from '../../../transactions/components/TransactionStateMessage.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'
import { useActiveAccounts } from '../../../transactions/hooks/useActiveAccounts.ts'
import { useActiveCategories } from '../../../transactions/hooks/useActiveCategories.ts'
import {
  formatMonthQuery,
  getCurrentSelectedMonth,
  getMonthPeriod,
} from '../../../transactions/utils/date.ts'
import { useCommitments } from '../../hooks/useCommitments.ts'
import { calculateCommitmentTotals } from '../../lib/commitmentTotals.ts'
import { parseMonthQuery } from '../../lib/commitmentDates.ts'
import {
  categoryTypeForCommitment,
  parseCommitmentType,
} from '../../lib/commitmentType.ts'
import type { CommitmentType } from '../../types/commitments.ts'
import { CommitmentList } from './CommitmentList.tsx'
import { PlanningHeader } from './PlanningHeader.tsx'
import { PlanningTabs } from './PlanningTabs.tsx'
import { PlanningTotals } from './PlanningTotals.tsx'

export function PlanningPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const [type, setType] = useState<CommitmentType>(parseCommitmentType(searchParams.get('type')))
  const selectedMonth = useMemo(
    () => parseMonthQuery(searchParams.get('month')) ?? getCurrentSelectedMonth(),
    [searchParams],
  )
  const period = useMemo(() => getMonthPeriod(selectedMonth), [selectedMonth])
  const commitmentsQuery = useCommitments(type, period)
  const accountsQuery = useActiveAccounts()
  const categoriesQuery = useActiveCategories(categoryTypeForCommitment(type))
  const accountNames = useMemo(() => mapById(accountsQuery.data), [accountsQuery.data])
  const categoryNames = useMemo(() => mapById(categoriesQuery.data), [categoriesQuery.data])
  const totals = useMemo(
    () => calculateCommitmentTotals(commitmentsQuery.data),
    [commitmentsQuery.data],
  )

  const updateSelectedMonth = useCallback(
    (nextMonth: SelectedMonth) => {
      const nextParams = new URLSearchParams(searchParams)
      nextParams.set('month', formatMonthQuery(nextMonth))
      nextParams.set('type', type)
      setSearchParams(nextParams)
    },
    [searchParams, setSearchParams, type],
  )

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

  function handleSelectType(nextType: CommitmentType) {
    setType(nextType)
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set('type', nextType)
    nextParams.set('month', formatMonthQuery(selectedMonth))
    setSearchParams(nextParams)
  }

  return (
    <TransactionsPageLayout animationKey={`${type}-${formatMonthQuery(selectedMonth)}`}>
      <section className="mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-hidden bg-[#6818e8] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none">
        <PlanningHeader
          selectedMonth={selectedMonth}
          onBack={() => navigate('/')}
          onCreate={() =>
            navigate(`/planning/${type === 'payable' ? 'payables' : 'receivables'}/new`)
          }
          onChangeMonth={changeSelectedMonth}
          onSelectMonth={updateSelectedMonth}
        />

        <div className="scrollbar-none flex min-h-0 flex-1 flex-col overflow-y-auto overflow-x-hidden rounded-t-[26px] bg-white px-5 pb-[var(--app-mobile-content-bottom)] pt-4 md:px-7 md:pb-10">
          <div className="flex w-full min-w-0 flex-1 flex-col gap-3">
            <PlanningTabs type={type} onSelectType={handleSelectType} />
            <PlanningTotals totals={totals} />

            {commitmentsQuery.isLoading ? (
              <TransactionStateMessage>Carregando compromissos...</TransactionStateMessage>
            ) : null}
            {commitmentsQuery.isError ? (
              <TransactionStateMessage tone="danger">
                Nao foi possivel carregar os compromissos.
              </TransactionStateMessage>
            ) : null}
            {commitmentsQuery.data ? (
              <CommitmentList
                type={type}
                commitments={commitmentsQuery.data}
                accountNames={accountNames}
                categoryNames={categoryNames}
              />
            ) : null}
          </div>
        </div>
      </section>
    </TransactionsPageLayout>
  )
}
