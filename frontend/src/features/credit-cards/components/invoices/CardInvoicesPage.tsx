import { ReceiptText } from 'lucide-react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { useCardInvoices } from '../../hooks/useCardInvoices.ts'
import { useCreditCards } from '../../hooks/useCreditCards.ts'
import { formatCardDate, formatInvoiceMonth } from '../../lib/creditCardDates.ts'
import { findById } from '../../lib/creditCardCollections.ts'
import { PageShell } from '../shared/PageShell.tsx'
import { StatePanel } from '../shared/PageState.tsx'
import { InvoiceBadge } from './InvoiceBadge.tsx'

export function CardInvoicesPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const cardId = searchParams.get('cardId') ?? ''
  const cardsQuery = useCreditCards()
  const invoicesQuery = useCardInvoices(cardId)
  const card = findById(cardsQuery.data, cardId)

  return (
    <PageShell title={card?.name ?? 'Faturas'}>
      {!cardId ? <StatePanel tone="danger">Cartao nao informado.</StatePanel> : null}
      {invoicesQuery.isLoading ? <StatePanel>Carregando faturas...</StatePanel> : null}
      {invoicesQuery.isError ? <StatePanel tone="danger">Nao foi possivel carregar as faturas.</StatePanel> : null}
      {invoicesQuery.data?.length === 0 ? <StatePanel>Ainda nao ha faturas para este cartao.</StatePanel> : null}
      {invoicesQuery.data?.length ? (
        <ul className="divide-y divide-[#f0ebf6]">
          {invoicesQuery.data.map((invoice) => (
            <li key={invoice.id} className="grid grid-cols-[40px_minmax(0,1fr)_auto] items-center gap-3 px-1 py-3">
              <span className="grid h-10 w-10 place-items-center rounded-full bg-[#fff8e8] text-[#9b6a12]">
                <ReceiptText className="h-4.5 w-4.5" aria-hidden="true" />
              </span>
              <div className="min-w-0">
                <h3 className="truncate text-[14px] font-semibold text-[#2c2237]">{formatInvoiceMonth(invoice.referenceMonth)}</h3>
                <p className="mt-1 truncate text-[12px] font-semibold text-[#81788c]">Vence {formatCardDate(invoice.dueAt)}</p>
              </div>
              <button type="button" className="text-right" onClick={() => navigate(`/credit-cards/invoice?invoiceId=${encodeURIComponent(invoice.id)}`)}>
                <strong className="block text-[14px] font-semibold text-[#c72f4d]">{formatCurrency(invoice.amount)}</strong>
                <span className="mt-1 block"><InvoiceBadge invoice={invoice} /></span>
              </button>
            </li>
          ))}
        </ul>
      ) : null}
    </PageShell>
  )
}
