import { CreditCard as CreditCardIcon, Plus, ReceiptText } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useConfirmDialog } from '../../../../components/confirm-dialog-context.ts'
import { ItemActionsMenu } from '../../../../components/ItemActionsMenu.tsx'
import type { Account } from '../../../transactions/types/transactions.ts'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { useCardInvoices } from '../../hooks/useCardInvoices.ts'
import { useInactivateCreditCard } from '../../hooks/useCreditCardMutations.domain.ts'
import { formatCardDate } from '../../lib/creditCardDates.ts'
import { cardStatusCopy } from '../../lib/creditCardPresentation.ts'
import type { CreditCard } from '../../types/credit-card.types.ts'
import { InvoiceBadge } from '../invoices/InvoiceBadge.tsx'

export function CreditCardRow({ card, accounts }: { card: CreditCard; accounts: Map<string, Account> }) {
  const navigate = useNavigate()
  const { confirm } = useConfirmDialog()
  const inactivateMutation = useInactivateCreditCard()
  const invoicesQuery = useCardInvoices(card.id)
  const linkedAccount = accounts.get(card.linkedAccountId)
  const currentInvoice = invoicesQuery.data?.find((invoice) => invoice.effectiveStatus === 'open')
    ?? invoicesQuery.data?.find((invoice) => invoice.effectiveStatus === 'closed' || invoice.effectiveStatus === 'overdue')

  async function handleInactivate() {
    const shouldInactivate = await confirm({
      title: 'Inativar cartao',
      description: `Inativar o cartao "${card.name}"?`,
      confirmLabel: 'Inativar',
      cancelLabel: 'Voltar',
      tone: 'danger',
    })

    if (shouldInactivate) {
      inactivateMutation.mutate(card.id)
    }
  }

  return (
    <li className="grid gap-3 px-1 py-3 sm:px-2">
      <div className="grid grid-cols-[40px_minmax(0,1fr)_minmax(88px,auto)_32px] items-center gap-3">
        <span className="grid h-10 w-10 place-items-center rounded-full bg-[#eef6ff] text-[#216fb8]">
          <CreditCardIcon className="h-4.5 w-4.5" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h3 className="truncate text-[14px] font-semibold leading-tight text-[#2c2237]">{card.name}</h3>
          <p className="mt-1 truncate text-[12px] font-semibold leading-tight text-[#81788c]">
            {linkedAccount?.name ?? 'Conta vinculada'} / fecha dia {card.closingDay}
          </p>
        </div>
        <div className="min-w-0 text-right">
          <strong className="block text-[13px] font-semibold leading-tight text-[#c72f4d] sm:text-[14px]">
            {formatCurrency(card.limitUsed)}
          </strong>
          <span className="mt-1 block text-[12px] font-semibold text-[#81788c]">
            livre {formatCurrency(card.limitAvailable)}
          </span>
        </div>
        <ItemActionsMenu
          label={`Acoes de ${card.name}`}
          onEdit={() => navigate(`/credit-cards/edit?cardId=${encodeURIComponent(card.id)}`)}
          onDelete={handleInactivate}
          isDeleteDisabled={card.status !== 'active' || inactivateMutation.isPending}
        />
      </div>

      <div className="grid gap-2 pl-[52px] text-[12px] font-semibold text-[#81788c]">
        <div className="flex min-w-0 flex-wrap items-center gap-2">
          <span className={card.status === 'active' ? 'text-[#147a46]' : 'text-[#81788c]'}>
            {cardStatusCopy[card.status]}
          </span>
          <span>limite {formatCurrency(card.limitTotal)}</span>
          <span>vence dia {card.dueDay}</span>
        </div>
        {currentInvoice ? (
          <div className="flex min-w-0 flex-wrap items-center gap-2">
            <InvoiceBadge invoice={currentInvoice} />
            <span>fatura {formatCurrency(currentInvoice.amount)}</span>
            <span>{formatCardDate(currentInvoice.dueAt)}</span>
          </div>
        ) : null}
      </div>

      <div className="flex flex-wrap gap-2 pl-[52px]">
        <button
          type="button"
          className="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-full bg-[#281d35] px-3 text-[12px] font-semibold text-white transition-colors hover:bg-[#3a2a4a]"
          onClick={() => navigate(`/credit-cards/purchase?cardId=${encodeURIComponent(card.id)}`)}
        >
          <Plus className="h-3.5 w-3.5" aria-hidden="true" />
          Compra
        </button>
        <button
          type="button"
          className="inline-flex h-9 cursor-pointer items-center gap-1.5 rounded-full border border-[#e3ddea] bg-white px-3 text-[12px] font-semibold text-[#4f435c] transition-colors hover:bg-[#f8f5fb]"
          onClick={() => navigate(`/credit-cards/invoices?cardId=${encodeURIComponent(card.id)}`)}
        >
          <ReceiptText className="h-3.5 w-3.5" aria-hidden="true" />
          Faturas
        </button>
      </div>
    </li>
  )
}
