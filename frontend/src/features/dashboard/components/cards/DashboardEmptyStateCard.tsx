import type { ReactNode } from 'react'

type DashboardEmptyStateCardProps = {
  title: string
  message: string
  action?: ReactNode
  className?: string
}

export function DashboardEmptyStateCard({
  title,
  message,
  action,
  className = '',
}: DashboardEmptyStateCardProps) {
  return (
    <div
      className={`grid place-items-center rounded-[18px] border border-[#ece8f2] bg-white px-5 py-8 text-center shadow-[0_16px_38px_rgba(48,39,61,0.07)] ${className}`}
    >
      <div className="grid max-w-[360px] justify-items-center gap-3">
        <div className="grid gap-1.5">
          <h3 className="m-0 text-[15px] font-semibold leading-tight text-[#241a30]">
            {title}
          </h3>
          <p className="m-0 text-[13px] font-medium leading-relaxed text-[#81798b]">
            {message}
          </p>
        </div>
        {action ? <div className="pt-1">{action}</div> : null}
      </div>
    </div>
  )
}
