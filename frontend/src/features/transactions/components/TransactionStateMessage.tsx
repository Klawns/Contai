import type { ReactNode } from 'react'

type TransactionStateMessageProps = {
  tone?: 'neutral' | 'danger'
  children: ReactNode
}

export function TransactionStateMessage({
  tone = 'neutral',
  children,
}: TransactionStateMessageProps) {
  const className =
    tone === 'danger'
      ? 'border-[#f0caca] text-[#b93838]'
      : 'border-[#ece8f2] text-[#81798b]'

  return (
    <div className={`rounded-[18px] border bg-[#fbf9fd] px-5 py-6 text-[14px] font-medium ${className}`}>
      {children}
    </div>
  )
}
