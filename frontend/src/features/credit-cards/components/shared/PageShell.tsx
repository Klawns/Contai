import { ArrowLeft } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { TransactionsPageLayout } from '../../../transactions/components/TransactionsPageLayout.tsx'

export function PageShell({
  title,
  children,
  action,
}: {
  title: string
  children: React.ReactNode
  action?: React.ReactNode
}) {
  const navigate = useNavigate()

  return (
    <TransactionsPageLayout animationKey={title}>
      <section className="mx-auto flex h-full min-h-0 w-full max-w-[520px] flex-col overflow-hidden bg-[#6818e8] text-left shadow-[0_24px_70px_rgba(43,35,54,0.12)] md:mx-0 md:max-w-none md:shadow-none">
        <header className="flex-none bg-[#6818e8] px-5 pb-5 pt-[calc(18px+env(safe-area-inset-top))] text-white md:px-7 md:pt-6">
          <div className="mx-auto grid w-full grid-cols-[44px_minmax(0,1fr)_44px] items-center">
            <button
              type="button"
              className="grid h-11 w-11 cursor-pointer place-items-center rounded-full bg-white/14 text-white transition-colors hover:bg-white/22 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              aria-label="Voltar"
              onClick={() => navigate('/more')}
            >
              <ArrowLeft className="h-5 w-5" aria-hidden="true" />
            </button>
            <h1 className="truncate px-2 text-center text-[17px] font-semibold leading-tight md:text-[24px]">
              {title}
            </h1>
            {action ?? <div aria-hidden="true" />}
          </div>
        </header>
        <div className="scrollbar-none flex min-h-0 flex-1 flex-col overflow-y-auto overflow-x-hidden rounded-t-[26px] bg-white px-5 pb-[var(--app-mobile-content-bottom)] pt-4 md:px-7 md:pb-10">
          <div className="flex w-full min-w-0 flex-1 flex-col gap-3">{children}</div>
        </div>
      </section>
    </TransactionsPageLayout>
  )
}
