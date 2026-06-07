import { invoiceStatusCopy } from '../../lib/invoicePresentation.ts'
import type { CardInvoice } from '../../types/invoice.types.ts'

export function InvoiceBadge({ invoice }: { invoice: CardInvoice }) {
  const tone =
    invoice.effectiveStatus === 'overdue'
      ? 'bg-[#fff0f2] text-[#c72f4d]'
      : invoice.effectiveStatus === 'paid'
        ? 'bg-[#e8f8ef] text-[#147a46]'
        : invoice.effectiveStatus === 'closed'
          ? 'bg-[#eef6ff] text-[#216fb8]'
          : invoice.effectiveStatus === 'canceled'
            ? 'bg-[#f1eef5] text-[#81788c]'
            : 'bg-[#fff8e8] text-[#9b6a12]'

  return (
    <span className={`rounded-full px-2.5 py-1 text-[11px] font-semibold ${tone}`}>
      {invoiceStatusCopy[invoice.effectiveStatus]}
    </span>
  )
}
