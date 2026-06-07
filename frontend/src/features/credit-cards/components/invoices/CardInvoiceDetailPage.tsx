import { useMemo } from 'react'
import { ReceiptText } from 'lucide-react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { mapById } from '../../../../lib/collections/mapById.ts'
import { useConfirmDialog } from '../../../../components/confirm-dialog-context.ts'
import { useActiveCategories } from '../../../transactions/hooks/useActiveCategories.ts'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { useCardInvoice } from '../../hooks/useCardInvoices.ts'
import { useCardPurchases } from '../../hooks/useCardPurchases.ts'
import { useCloseCardInvoice } from '../../hooks/useCardInvoiceMutations.ts'
import { useCreditCards } from '../../hooks/useCreditCards.ts'
import { formatCardDate, formatInvoiceMonth } from '../../lib/creditCardDates.ts'
import { findById } from '../../lib/creditCardCollections.ts'
import { PageShell } from '../shared/PageShell.tsx'
import { StatePanel } from '../shared/PageState.tsx'
import { InvoiceBadge } from './InvoiceBadge.tsx'
import { InvoiceInstallmentsList } from './InvoiceInstallmentsList.tsx'

export function CardInvoiceDetailPage() {
  const navigate = useNavigate()
  const { confirm } = useConfirmDialog()
  const [searchParams] = useSearchParams()
  const invoiceId = searchParams.get('invoiceId') ?? ''
  const invoiceQuery = useCardInvoice(invoiceId)
  const cardsQuery = useCreditCards()
  const purchasesQuery = useCardPurchases(invoiceQuery.data?.cardId ?? '')
  const categoriesQuery = useActiveCategories('expense')
  const closeMutation = useCloseCardInvoice()
  const invoice = invoiceQuery.data
  const card = findById(cardsQuery.data, invoice?.cardId ?? '')
  const purchaseNames = useMemo(() => mapById(purchasesQuery.data), [purchasesQuery.data])
  const categoryNames = useMemo(() => mapById(categoriesQuery.data), [categoriesQuery.data])

  async function handleClose() {
    const shouldClose = await confirm({
      title: 'Fechar fatura',
      description: 'Fechar esta fatura para pagamento?',
      confirmLabel: 'Fechar',
      cancelLabel: 'Voltar',
      tone: 'default',
    })

    if (shouldClose) {
      closeMutation.mutate(invoiceId)
    }
  }

  return (
    <PageShell title={card?.name ?? 'Fatura'}>
      {!invoiceId ? <StatePanel tone="danger">Fatura nao informada.</StatePanel> : null}
      {invoiceQuery.isLoading ? <StatePanel>Carregando fatura...</StatePanel> : null}
      {invoiceQuery.isError ? <StatePanel tone="danger">Nao foi possivel carregar a fatura.</StatePanel> : null}
      {invoice ? (
        <>
          <section className="grid gap-3 border-b border-[#f0ebf6] pb-4">
            <div>
              <ReceiptText className="h-5 w-5 text-[#9b6a12]" aria-hidden="true" />
              <span className="mt-2 block text-[12px] font-semibold text-[#81788c]">{formatInvoiceMonth(invoice.referenceMonth)}</span>
              <strong className="mt-1 block text-[20px] font-semibold text-[#c72f4d]">{formatCurrency(invoice.amount)}</strong>
            </div>
            <div className="flex flex-wrap gap-2 text-[12px] font-semibold text-[#81788c]">
              <InvoiceBadge invoice={invoice} />
              <span>Fecha {formatCardDate(invoice.closingAt)}</span>
              <span>Vence {formatCardDate(invoice.dueAt)}</span>
            </div>
            <div className="flex flex-wrap gap-2">
              {invoice.effectiveStatus === 'open' ? (
                <button type="button" disabled={closeMutation.isPending} className="h-9 rounded-full bg-[#281d35] px-3 text-[12px] font-semibold text-white disabled:opacity-60" onClick={() => void handleClose()}>
                  Fechar fatura
                </button>
              ) : null}
              {invoice.effectiveStatus === 'closed' || invoice.effectiveStatus === 'overdue' ? (
                <button type="button" className="h-9 rounded-full bg-[#147a46] px-3 text-[12px] font-semibold text-white" onClick={() => navigate(`/credit-cards/invoice/pay?invoiceId=${encodeURIComponent(invoice.id)}`)}>
                  Pagar fatura
                </button>
              ) : null}
            </div>
          </section>
          <InvoiceInstallmentsList
            invoice={invoice}
            purchases={purchaseNames}
            categories={categoryNames}
          />
        </>
      ) : null}
    </PageShell>
  )
}
