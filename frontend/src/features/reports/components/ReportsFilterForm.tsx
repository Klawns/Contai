import { useMemo } from 'react'
import type { FormEvent } from 'react'
import {
  CalendarDays,
  Download,
  FolderTree,
  Landmark,
  Layers3,
  Tag,
} from 'lucide-react'
import {
  AccountOption,
  CategoryOption,
  TransactionSelectField,
} from '../../transactions/components/Selectors.tsx'
import { FormScrollableContent } from '../../transactions/components/FormFields.tsx'
import type { TransactionSelectOption } from '../../transactions/components/Selectors.tsx'
import type { Account, Category } from '../../transactions/types/transactions.ts'
import {
  groupOptions,
  movementTypeOptions,
  settlementOptions,
} from '../data/reportFilterOptions.ts'
import { useReportAccounts } from '../hooks/useReportAccounts.ts'
import { useReportCategories } from '../hooks/useReportCategories.ts'
import type { ReportFormState } from '../types/reportForm.ts'
import { allReportOption } from '../utils/reportFilters.ts'
import { ReportAllOptionButton } from './ReportAllOptionButton.tsx'
import { ReportDateField } from './ReportDateField.tsx'

type ReportsFilterFormProps = {
  formState: ReportFormState
  isExporting: boolean
  message: string
  onFieldChange: <Key extends keyof ReportFormState>(
    key: Key,
    value: ReportFormState[Key],
  ) => void
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}

export function ReportsFilterForm({
  formState,
  isExporting,
  message,
  onFieldChange,
  onSubmit,
}: ReportsFilterFormProps) {
  const accountsQuery = useReportAccounts()
  const categoriesQuery = useReportCategories()
  const accountOptions = useMemo<Array<TransactionSelectOption<string, Account>>>(
    () => [
      { value: allReportOption, label: 'Todas' },
      ...(accountsQuery.data ?? []).map((account) => ({
        value: account.id,
        label: account.name,
        description: account.type,
        item: account,
      })),
    ],
    [accountsQuery.data],
  )
  const categoryOptions = useMemo<Array<TransactionSelectOption<string, Category>>>(
    () => [
      { value: allReportOption, label: 'Todas' },
      ...(categoriesQuery.data ?? []).map((category) => ({
        value: category.id,
        label: category.name,
        description: category.type === 'income' ? 'Receita' : 'Despesa',
        item: category,
      })),
    ],
    [categoriesQuery.data],
  )
  const selectedCategory = categoryOptions.find((option) => option.value === formState.categoryId)

  return (
    <form
      className="scrollbar-none flex min-h-0 w-full flex-1 flex-col overflow-y-auto overflow-x-hidden rounded-t-[28px] bg-white shadow-[0_-1px_8px_rgba(17,24,39,0.04)]"
      onSubmit={onSubmit}
    >
      <FormScrollableContent className="pt-2">
        <ReportDateField
          label="Data inicial"
          value={formState.startDate}
          accentColor="#7b2cff"
          onChange={(value) => onFieldChange('startDate', value)}
        />
        <ReportDateField
          label="Data final"
          value={formState.endDate}
          accentColor="#c72f4d"
          onChange={(value) => onFieldChange('endDate', value)}
        />
        <TransactionSelectField
          label="Tipo"
          value={formState.movementType}
          placeholder="Todos"
          icon={<Layers3 className="h-5 w-5" aria-hidden="true" />}
          options={movementTypeOptions}
          sheetTitle="Tipo de movimentacao"
          onChange={(value) => onFieldChange('movementType', value)}
        />
        <TransactionSelectField
          label="Categoria"
          value={formState.categoryId}
          placeholder="Todas"
          icon={<Tag className="h-5 w-5" aria-hidden="true" />}
          options={categoryOptions}
          sheetTitle="Categoria"
          onChange={(value) => onFieldChange('categoryId', value)}
          chipClassName={selectedCategory?.item ? 'text-white' : undefined}
          chipStyle={selectedCategory?.item ? { backgroundColor: selectedCategory.item.color } : undefined}
          isLoading={categoriesQuery.isLoading}
          loadingMessage="Carregando categorias..."
          isError={categoriesQuery.isError}
          errorMessage="Nao foi possivel carregar as categorias."
          emptyMessage="Nenhuma categoria disponivel."
          renderOption={({ option, isSelected, onSelect }) => (
            option.item ? (
              <CategoryOption
                category={option.item}
                isSelected={isSelected}
                onSelect={onSelect}
              />
            ) : (
              <ReportAllOptionButton label={option.label} isSelected={isSelected} onSelect={onSelect} />
            )
          )}
        />
        <TransactionSelectField
          label="Conta"
          value={formState.accountId}
          placeholder="Todas"
          icon={<Landmark className="h-5 w-5" aria-hidden="true" />}
          options={accountOptions}
          sheetTitle="Conta"
          onChange={(value) => onFieldChange('accountId', value)}
          isLoading={accountsQuery.isLoading}
          loadingMessage="Carregando contas..."
          isError={accountsQuery.isError}
          errorMessage="Nao foi possivel carregar as contas."
          emptyMessage="Nenhuma conta disponivel."
          renderOption={({ option, isSelected, onSelect }) => (
            option.item ? (
              <AccountOption
                account={option.item}
                isSelected={isSelected}
                onSelect={onSelect}
              />
            ) : (
              <ReportAllOptionButton label={option.label} isSelected={isSelected} onSelect={onSelect} />
            )
          )}
        />
        <TransactionSelectField
          label="Status"
          value={formState.settlementStatus}
          placeholder="Todos"
          icon={<CalendarDays className="h-5 w-5" aria-hidden="true" />}
          options={settlementOptions}
          sheetTitle="Status financeiro"
          onChange={(value) => onFieldChange('settlementStatus', value)}
        />
        <TransactionSelectField
          label="Agrupamento"
          value={formState.groupBy}
          placeholder="Nenhum"
          icon={<FolderTree className="h-5 w-5" aria-hidden="true" />}
          options={groupOptions}
          sheetTitle="Agrupamento"
          onChange={(value) => onFieldChange('groupBy', value)}
        />

        {message ? (
          <p className="mx-4 mt-4 rounded-lg border border-[#f1c8d0] bg-[#fff7f9] px-4 py-3 text-[14px] font-semibold text-[#c72f4d] md:mx-5">
            {message}
          </p>
        ) : null}
      </FormScrollableContent>

      <div className="sticky bottom-[var(--app-mobile-sticky-bottom)] bg-white/96 px-4 pb-4 pt-5 backdrop-blur md:static md:bg-transparent md:px-5 md:pb-5 md:pt-6 md:backdrop-blur-none">
        <button
          type="submit"
          disabled={isExporting}
          className="mx-auto flex h-12 w-full max-w-[420px] cursor-pointer items-center justify-center gap-2 rounded-lg bg-[#281d35] px-4 text-[15px] font-semibold text-white shadow-[0_6px_14px_rgba(40,29,53,0.10)] transition-colors hover:bg-[#3a2a4a] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] disabled:cursor-not-allowed disabled:opacity-65"
        >
          <Download className="h-4 w-4" aria-hidden="true" />
          {isExporting ? 'Gerando...' : 'Gerar relatorio'}
        </button>
      </div>
    </form>
  )
}
