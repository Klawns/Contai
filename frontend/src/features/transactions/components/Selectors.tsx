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

function SelectorButton({
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
  const [isOpen, setIsOpen] = useState(false)
  const accountsQuery = useActiveAccounts()
  const selected = accountsQuery.data?.find((account) => account.id === value)

  return (
    <>
      <SelectorButton
        label={label}
        valueLabel={selected?.name ?? 'Selecione uma conta'}
        error={error}
        icon={<Landmark className="h-5 w-5" aria-hidden="true" />}
        chipClassName={selected ? 'bg-[#eef6ff] text-[#216fb8]' : undefined}
        onClick={() => setIsOpen(true)}
      />
      <SelectionSheet title={label} isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="grid gap-2">
          {accountsQuery.isLoading ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#81788c]">Carregando contas...</p>
          ) : null}
          {accountsQuery.isError ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#b93838]">
              Nao foi possivel carregar as contas.
            </p>
          ) : null}
          {(accountsQuery.data ?? []).map((account) => (
            <AccountOption
              key={account.id}
              account={account}
              isSelected={account.id === value}
              onSelect={() => {
                onChange(account.id)
                setIsOpen(false)
              }}
            />
          ))}
        </div>
      </SelectionSheet>
    </>
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
  const [isOpen, setIsOpen] = useState(false)
  const categoriesQuery = useActiveCategories(type)
  const selected = useMemo(
    () => categoriesQuery.data?.find((category) => category.id === value),
    [categoriesQuery.data, value],
  )

  return (
    <>
      <SelectorButton
        label="Categoria"
        valueLabel={selected?.name ?? 'Selecione uma categoria'}
        error={error}
        icon={<Tag className="h-5 w-5" aria-hidden="true" />}
        chipClassName={selected ? 'text-white' : undefined}
        chipStyle={selected ? { backgroundColor: selected.color } : undefined}
        onClick={() => setIsOpen(true)}
      />
      <SelectionSheet title="Categoria" isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="grid gap-2">
          <AddCategoryButton
            onClick={() => {
              setIsOpen(false)
              onAddCategory()
            }}
          />
          {categoriesQuery.isLoading ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#81788c]">
              Carregando categorias...
            </p>
          ) : null}
          {categoriesQuery.isError ? (
            <p className="px-1 py-2 text-[14px] font-medium text-[#b93838]">
              Nao foi possivel carregar as categorias.
            </p>
          ) : null}
          {(categoriesQuery.data ?? []).map((category) => (
            <CategoryOption
              key={category.id}
              category={category}
              isSelected={category.id === value}
              onSelect={() => {
                onChange(category.id)
                setIsOpen(false)
              }}
            />
          ))}
        </div>
      </SelectionSheet>
    </>
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
