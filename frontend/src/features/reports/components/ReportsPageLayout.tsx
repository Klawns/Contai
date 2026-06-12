import type { ReactNode } from 'react'

type ReportsPageLayoutProps = {
  children: ReactNode
}

export function ReportsPageLayout({ children }: ReportsPageLayoutProps) {
  return (
    <section className="flex h-full min-h-0 w-full max-w-none flex-col overflow-hidden bg-[#eaf3fb]">
      <header className="grid w-full flex-none gap-4 px-4 pb-4 pt-[calc(16px+env(safe-area-inset-top))] sm:px-5 md:px-8 md:pt-6 lg:px-10">
        <h1 className="text-center text-[18px] font-semibold leading-tight text-[#18202f]">
          Relatorios
        </h1>
      </header>

      {children}
    </section>
  )
}
