import { useMemo } from 'react'
import { CreditCard as CreditCardIcon } from 'lucide-react'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { TransactionSelectField } from '../../../transactions/components/Selectors.tsx'
import { useCreditCards } from '../../hooks/useCreditCards.ts'
import { cardStatusOptions } from '../../lib/creditCardPresentation.ts'
import type { CreditCard, CreditCardStatus } from '../../types/credit-card.types.ts'

export function CardSelector({
  value,
  error,
  onChange,
}: {
  value: string
  error?: string
  onChange: (value: string) => void
}) {
  const cardsQuery = useCreditCards()
  const options = useMemo(
    () =>
      (cardsQuery.data ?? [])
        .filter((card) => card.status === 'active')
        .map((card) => ({
          value: card.id,
          label: card.name,
          description: `Disponivel ${formatCurrency(card.limitAvailable)}`,
          item: card,
        })),
    [cardsQuery.data],
  )
  const selected = options.find((option) => option.value === value)

  return (
    <TransactionSelectField
      label="Cartao"
      value={value}
      placeholder="Selecione um cartao"
      error={error}
      icon={<CreditCardIcon className="h-5 w-5" aria-hidden="true" />}
      options={options}
      onChange={onChange}
      chipClassName={selected ? 'bg-[#eef6ff] text-[#216fb8]' : undefined}
      isLoading={cardsQuery.isLoading}
      loadingMessage="Carregando cartoes..."
      isError={cardsQuery.isError}
      errorMessage="Nao foi possivel carregar os cartoes."
      emptyMessage="Nenhum cartao ativo disponivel."
      renderOption={({ option, isSelected, onSelect }) => (
        option.item ? (
          <CardOption card={option.item} isSelected={isSelected} onSelect={onSelect} />
        ) : null
      )}
    />
  )
}

export function CardStatusSelector({
  value,
  error,
  onChange,
}: {
  value: CreditCardStatus
  error?: string
  onChange: (value: CreditCardStatus) => void
}) {
  const selected = cardStatusOptions.find((option) => option.value === value)

  return (
    <TransactionSelectField
      label="Status"
      value={value}
      placeholder="Selecione o status"
      error={error}
      icon={<CreditCardIcon className="h-5 w-5" aria-hidden="true" />}
      options={cardStatusOptions}
      onChange={onChange}
      chipClassName={selected?.value === 'active' ? 'bg-[#e8f8ef] text-[#147a46]' : 'bg-[#f1eef5] text-[#81788c]'}
      sheetTitle="Status"
    />
  )
}

function CardOption({
  card,
  isSelected,
  onSelect,
}: {
  card: CreditCard
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
      <span className="grid h-10 w-10 flex-none place-items-center rounded-full bg-[#eef6ff] text-[#216fb8]">
        <CreditCardIcon className="h-5 w-5" aria-hidden="true" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="block truncate text-[14px] font-semibold text-[#2c2237]">
          {card.name}
        </span>
        <span className="block truncate text-[12px] font-medium text-[#81788c]">
          Disponivel {formatCurrency(card.limitAvailable)}
        </span>
      </span>
    </button>
  )
}
