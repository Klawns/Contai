import { useMemo, useState } from 'react'
import {
  BriefcaseBusiness,
  Car,
  ChevronRight,
  CirclePlus,
  Ellipsis,
  GraduationCap,
  HeartPulse,
  House,
  Landmark,
  Laptop,
  PartyPopper,
  Plus,
  ShoppingBag,
  Tag,
  TrendingUp,
  Utensils,
} from 'lucide-react'
import type { CSSProperties, ReactNode } from 'react'
import type { LucideIcon } from 'lucide-react'
import { useActiveAccounts } from '../hooks/useActiveAccounts.ts'
import { useActiveCategories } from '../hooks/useActiveCategories.ts'
import type { Account, Category, CategoryTransactionType } from '../types/transactions.ts'
import { formatCurrency } from '../utils/money.ts'
import { TransactionFieldRow } from './FormFields.tsx'
import { SelectionSheet } from './SelectionSheet.tsx'

type SelectorButtonProps = {
  label: string
  valueLabel: string
  error?: string
  icon: ReactNode
  chipClassName?: string
  chipStyle?: CSSProperties
  onClick: () => void
}

function TransactionSelectTrigger({
  label,
  valueLabel,
  error,
  icon,
  chipClassName = 'bg-[#f4f1f7] text-[#2c2237]',
  chipStyle,
  onClick,
}: SelectorButtonProps) {
  return (
    <TransactionFieldRow
      icon={icon}
      label={label}
      error={error}
    >
      <button
        type="button"
        className="flex min-h-11 w-full cursor-pointer items-center justify-between gap-2 text-left focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        onClick={onClick}
      >
        <span
          className={`min-w-0 max-w-full truncate rounded-full px-3 py-1.5 text-[13px] font-semibold ${chipClassName}`}
          style={chipStyle}
        >
          {valueLabel}
        </span>
        <ChevronRight className="h-4 w-4 flex-none text-[#9a91a5]" aria-hidden="true" />
      </button>
    </TransactionFieldRow>
  )
}

export type TransactionSelectOption<TValue extends string = string, TItem = unknown> = {
  value: TValue
  label: string
  description?: string
  item?: TItem
}

type TransactionSelectFieldProps<TValue extends string = string, TItem = unknown> = {
  label: string
  value: TValue | ''
  placeholder: string
  icon: ReactNode
  options: Array<TransactionSelectOption<TValue, TItem>>
  onChange: (value: TValue) => void
  error?: string
  sheetTitle?: string
  chipClassName?: string
  chipStyle?: CSSProperties
  loadingMessage?: string
  isLoading?: boolean
  errorMessage?: string
  isError?: boolean
  emptyMessage?: string
  beforeOptions?: (close: () => void) => ReactNode
  renderOption?: (params: {
    option: TransactionSelectOption<TValue, TItem>
    isSelected: boolean
    onSelect: () => void
  }) => ReactNode
}

export function TransactionSelectField<TValue extends string = string, TItem = unknown>({
  label,
  value,
  placeholder,
  icon,
  options,
  onChange,
  error,
  sheetTitle = label,
  chipClassName,
  chipStyle,
  loadingMessage = 'Carregando...',
  isLoading = false,
  errorMessage = 'Nao foi possivel carregar as opcoes.',
  isError = false,
  emptyMessage = 'Nenhuma opcao disponivel.',
  beforeOptions,
  renderOption,
}: TransactionSelectFieldProps<TValue, TItem>) {
  const [isOpen, setIsOpen] = useState(false)
  const selected = options.find((option) => option.value === value)
  const close = () => setIsOpen(false)

  function handleSelect(nextValue: TValue) {
    onChange(nextValue)
    close()
  }

  return (
    <>
      <TransactionSelectTrigger
        label={label}
        valueLabel={selected?.label ?? placeholder}
        error={error}
        icon={icon}
        chipClassName={chipClassName}
        chipStyle={chipStyle}
        onClick={() => setIsOpen(true)}
      />
      <SelectionSheet title={sheetTitle} isOpen={isOpen} onClose={close}>
        <div className="grid gap-2">
          {beforeOptions ? beforeOptions(close) : null}
          {isLoading ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#81788c]">{loadingMessage}</p>
          ) : null}
          {isError ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#b93838]">
              {errorMessage}
            </p>
          ) : null}
          {!isLoading && !isError && options.length === 0 ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#81788c]">{emptyMessage}</p>
          ) : null}
          {options.map((option) => {
            const isSelected = option.value === value
            const onSelect = () => handleSelect(option.value)

            return renderOption ? (
              <div key={option.value}>{renderOption({ option, isSelected, onSelect })}</div>
            ) : (
              <DefaultSelectOption
                key={option.value}
                option={option}
                isSelected={isSelected}
                onSelect={onSelect}
              />
            )
          })}
        </div>
      </SelectionSheet>
    </>
  )
}

function DefaultSelectOption<TValue extends string, TItem>({
  option,
  isSelected,
  onSelect,
}: {
  option: TransactionSelectOption<TValue, TItem>
  isSelected: boolean
  onSelect: () => void
}) {
  return (
    <button
      type="button"
      className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
        isSelected ? 'border-[#7b2cff] bg-[#f7f2ff]' : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
      }`}
      onClick={onSelect}
    >
      <span className="min-w-0 flex-1">
        <span className="block truncate text-[14px] font-semibold text-[#2c2237]">
          {option.label}
        </span>
        {option.description ? (
          <span className="block truncate text-[12px] font-medium text-[#81788c]">
            {option.description}
          </span>
        ) : null}
      </span>
    </button>
  )
}

type AccountOptionProps = {
  account: Account
  isSelected: boolean
  onSelect: () => void
}

export function AccountOption({ account, isSelected, onSelect }: AccountOptionProps) {
  return (
    <button
      type="button"
      className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
        isSelected ? 'border-[#7b2cff] bg-[#f7f2ff]' : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
      }`}
      onClick={onSelect}
    >
      <span className="grid h-10 w-10 flex-none place-items-center rounded-full bg-[#f0ebf8] text-[#6a22e5]">
        <Landmark className="h-5 w-5" aria-hidden="true" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="block truncate text-[14px] font-semibold text-[#2c2237]">
          {account.name}
        </span>
        <span className="block truncate text-[12px] font-medium text-[#81788c]">
          {formatCurrency(account.currentBalance)}
        </span>
      </span>
    </button>
  )
}

type AccountSelectorProps = {
  label: string
  value: string
  error?: string
  onChange: (value: string) => void
}

export function AccountSelector({ label, value, error, onChange }: AccountSelectorProps) {
  const accountsQuery = useActiveAccounts()
  const options = useMemo(
    () =>
      (accountsQuery.data ?? []).map((account) => ({
        value: account.id,
        label: account.name,
        description: formatCurrency(account.currentBalance),
        item: account,
      })),
    [accountsQuery.data],
  )
  const selected = options.find((option) => option.value === value)

  return (
    <TransactionSelectField
      label={label}
      value={value}
      placeholder="Selecione uma conta"
      error={error}
      icon={<Landmark className="h-5 w-5" aria-hidden="true" />}
      options={options}
      onChange={onChange}
      chipClassName={selected ? 'bg-[#eef6ff] text-[#216fb8]' : undefined}
      isLoading={accountsQuery.isLoading}
      loadingMessage="Carregando contas..."
      isError={accountsQuery.isError}
      errorMessage="Nao foi possivel carregar as contas."
      emptyMessage="Nenhuma conta disponivel."
      renderOption={({ option, isSelected, onSelect }) => (
        option.item ? (
          <AccountOption account={option.item} isSelected={isSelected} onSelect={onSelect} />
        ) : null
      )}
    />
  )
}

type CategoryOptionProps = {
  category: Category
  isSelected: boolean
  onSelect: () => void
}

const categoryIconById: Record<string, LucideIcon> = {
  laptop: Laptop,
  'trending-up': TrendingUp,
  'circle-plus': CirclePlus,
  'briefcase-business': BriefcaseBusiness,
  utensils: Utensils,
  'shopping-bag': ShoppingBag,
  'graduation-cap': GraduationCap,
  'party-popper': PartyPopper,
  house: House,
  ellipsis: Ellipsis,
  'heart-pulse': HeartPulse,
  car: Car,
  tag: Tag,
}

export function CategoryOption({ category, isSelected, onSelect }: CategoryOptionProps) {
  const Icon = categoryIconById[category.icon] ?? Tag

  return (
    <button
      type="button"
      className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
        isSelected ? 'border-[#7b2cff] bg-[#f7f2ff]' : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
      }`}
      onClick={onSelect}
    >
      <span
        className="grid h-10 w-10 flex-none place-items-center rounded-full text-white"
        style={{ backgroundColor: category.color }}
      >
        <Icon className="h-5 w-5" aria-hidden="true" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="block truncate text-[14px] font-semibold text-[#2c2237]">
          {category.name}
        </span>
      </span>
    </button>
  )
}

type CategorySelectorProps = {
  type: CategoryTransactionType
  value: string
  error?: string
  onChange: (value: string) => void
  onAddCategory: () => void
}

export function CategorySelector({
  type,
  value,
  error,
  onChange,
  onAddCategory,
}: CategorySelectorProps) {
  const categoriesQuery = useActiveCategories(type)
  const options = useMemo(
    () =>
      (categoriesQuery.data ?? []).map((category) => ({
        value: category.id,
        label: category.name,
        item: category,
      })),
    [categoriesQuery.data],
  )
  const selected = options.find((option) => option.value === value)

  return (
    <TransactionSelectField
      label="Categoria"
      value={value}
      placeholder="Selecione uma categoria"
      error={error}
      icon={<Tag className="h-5 w-5" aria-hidden="true" />}
      options={options}
      onChange={onChange}
      chipClassName={selected ? 'text-white' : undefined}
      chipStyle={selected?.item ? { backgroundColor: selected.item.color } : undefined}
      isLoading={categoriesQuery.isLoading}
      loadingMessage="Carregando categorias..."
      isError={categoriesQuery.isError}
      errorMessage="Nao foi possivel carregar as categorias."
      emptyMessage="Nenhuma categoria disponivel."
      beforeOptions={(close) => (
        <AddCategoryButton
          onClick={() => {
            close()
            onAddCategory()
          }}
        />
      )}
      renderOption={({ option, isSelected, onSelect }) => (
        option.item ? (
          <CategoryOption category={option.item} isSelected={isSelected} onSelect={onSelect} />
        ) : null
      )}
    />
  )
}

type AddCategoryButtonProps = {
  onClick: () => void
}

export function AddCategoryButton({ onClick }: AddCategoryButtonProps) {
  return (
    <button
      type="button"
      className="mb-2 flex h-11 w-full cursor-pointer items-center justify-center gap-2 rounded-lg border border-dashed border-[#cdbfe0] bg-[#fbf9fe] text-[14px] font-semibold text-[#6a22e5] transition-colors hover:bg-[#f5efff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
      onClick={onClick}
    >
      <Plus className="h-4 w-4" aria-hidden="true" />
      Adicionar categoria
    </button>
  )
}
