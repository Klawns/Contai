import { Tag, XCircle } from 'lucide-react'
import { useConfirmDialog } from '../../../../components/confirm-dialog-context.ts'
import type { Category } from '../../../transactions/types/transactions.ts'
import { formatCurrency } from '../../../transactions/utils/money.ts'
import { useCancelCardPurchase } from '../../hooks/useCardPurchaseMutations.ts'
import type { CardInvoice } from '../../types/invoice.types.ts'
import type { CardPurchase } from '../../types/purchase.types.ts'
import { StatePanel } from '../shared/PageState.tsx'

export function InvoiceInstallmentsList({
  invoice,
  purchases,
  categories,
}: {
  invoice: CardInvoice
  purchases: Map<string, CardPurchase>
  categories: Map<string, Category>
}) {
  const { confirm } = useConfirmDialog()
  const cancelMutation = useCancelCardPurchase()

  async function handleCancel(purchase: CardPurchase) {
    const shouldCancel = await confirm({
      title: 'Cancelar compra',
      description: `Cancelar "${purchase.description}"?`,
      confirmLabel: 'Cancelar compra',
      cancelLabel: 'Voltar',
      tone: 'danger',
    })

    if (shouldCancel) {
      cancelMutation.mutate(purchase.id)
    }
  }

  if (invoice.installments.length === 0) {
    return <StatePanel>Esta fatura ainda nao possui parcelas.</StatePanel>
  }

  return (
    <ul className="divide-y divide-[#f0ebf6]">
      {invoice.installments.map((installment) => {
        const purchase = purchases.get(installment.purchaseId)
        const category = purchase ? categories.get(purchase.categoryId) : undefined

        return (
          <li key={installment.id} className="grid grid-cols-[40px_minmax(0,1fr)_auto_32px] items-center gap-3 px-1 py-3">
            <span className="grid h-10 w-10 place-items-center rounded-full bg-[#fff0f2] text-[#c72f4d]">
              <Tag className="h-4.5 w-4.5" aria-hidden="true" />
            </span>
            <div className="min-w-0">
              <h3 className="truncate text-[14px] font-semibold text-[#2c2237]">{purchase?.description ?? 'Compra'}</h3>
              <p className="mt-1 truncate text-[12px] font-semibold text-[#81788c]">
                Parcela {installment.number}{category ? ` / ${category.name}` : ''}
              </p>
            </div>
            <strong className="text-[14px] font-semibold text-[#c72f4d]">{formatCurrency(installment.amount)}</strong>
            {purchase?.status === 'active' ? (
              <button
                type="button"
                className="grid h-8 w-8 place-items-center rounded-full text-[#c75959] hover:bg-[#fff4f4]"
                aria-label={`Cancelar ${purchase.description}`}
                onClick={() => void handleCancel(purchase)}
              >
                <XCircle className="h-4 w-4" aria-hidden="true" />
              </button>
            ) : (
              <span className="grid h-8 w-8 place-items-center text-[#b2a9bd]">-</span>
            )}
          </li>
        )
      })}
    </ul>
  )
}
