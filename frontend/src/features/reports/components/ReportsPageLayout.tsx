import type { ReactNode } from 'react'

type ReportsPageLayoutProps = {
  children: ReactNode
}

export function ReportsPageLayout({ children }: ReportsPageLayoutProps) {
  return (
    <section className="flex min-h-svh w-full max-w-none flex-col bg-[#eaf3fb] md:overflow-hidden">
      <header className="grid w-full gap-4 px-4 pb-4 pt-[calc(16px+env(safe-area-inset-top))] sm:px-5 md:px-8 md:pt-6 lg:px-10">
        <h1 className="text-center text-[18px] font-semibold leading-tight text-[#18202f]">
          Relatorios
        </h1>
      </header>

      {children}
    </section>
  )
}
