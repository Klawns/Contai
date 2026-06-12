import { TransactionStateMessage } from '../../../transactions/components/TransactionStateMessage.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'

type CommitmentPageStateProps = {
  animationKey: string
  tone: 'income' | 'expense'
  message: string
  danger?: boolean
  onBack?: () => void
}

export function CommitmentPageState({
  animationKey,
  tone,
  message,
  danger = false,
  onBack,
}: CommitmentPageStateProps) {
  return (
    <TransactionsPageLayout variant="create" tone={tone} animationKey={animationKey}>
      <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:px-8 md:shadow-none">
        {onBack ? (
          <button
            type="button"
            className="self-start text-[14px] font-semibold text-[#6818e8] transition-colors hover:text-[#4d12b0] focus-visible:rounded-md focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#6818e8]"
            onClick={onBack}
          >
            Voltar
          </button>
        ) : null}
        <TransactionStateMessage tone={danger ? 'danger' : undefined}>
          {message}
        </TransactionStateMessage>
      </section>
    </TransactionsPageLayout>
  )
}
