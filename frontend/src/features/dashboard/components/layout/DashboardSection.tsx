import type { ReactNode } from 'react'

type DashboardSectionProps = {
  title: string
  children: ReactNode
}

export function DashboardSection({ title, children }: DashboardSectionProps) {
  return (
    <section className="grid gap-2.5" aria-labelledby={`${title}-dashboard-section`}>
      <h2
        id={`${title}-dashboard-section`}
        className="m-0 px-1 text-[13px] font-medium leading-none text-[#81798b]"
      >
        {title}
      </h2>
      {children}
    </section>
  )
}
