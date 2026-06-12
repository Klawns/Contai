import { TransactionStateMessage } from '../../../transactions/components/TransactionStateMessage.tsx'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'

export function StatePanel({
  tone = 'default',
  children,
}: {
  tone?: 'default' | 'danger'
  children: React.ReactNode
}) {
  return (
    <TransactionStateMessage tone={tone === 'danger' ? 'danger' : 'neutral'}>
      {children}
    </TransactionStateMessage>
  )
}

export function FormLoading({ message }: { message: string }) {
  return (
    <TransactionsPageLayout variant="create" tone="expense" animationKey="form-loading">
      <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left md:mx-0 md:max-w-none md:px-8">
        <TransactionStateMessage>{message}</TransactionStateMessage>
      </section>
    </TransactionsPageLayout>
  )
}

export function FormInvalid({ message, onBack }: { message: string; onBack: () => void }) {
  return (
    <TransactionsPageLayout variant="create" tone="expense" animationKey="form-invalid">
      <section className="scrollbar-none mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-y-auto overflow-x-hidden bg-white px-5 py-[calc(28px+env(safe-area-inset-top))] text-left md:mx-0 md:max-w-none md:px-8">
        <button type="button" className="self-start text-[14px] font-semibold text-[#6818e8]" onClick={onBack}>Voltar</button>
        <TransactionStateMessage tone="danger">{message}</TransactionStateMessage>
      </section>
    </TransactionsPageLayout>
  )
}

export function ErrorMessage({ children }: { children: React.ReactNode }) {
  return (
    <p className="mx-4 mt-4 rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838] md:mx-5">
      {children}
    </p>
  )
}
