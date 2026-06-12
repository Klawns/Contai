import { CreditCard, FileText, Landmark } from 'lucide-react'
import { Link } from 'react-router-dom'
import { DashboardLayout } from '../dashboard/components'

const links = [
  { to: '/accounts', label: 'Contas', icon: Landmark },
  { to: '/credit-cards', label: 'Cartoes', icon: CreditCard },
  { to: '/reports', label: 'Relatorios', icon: FileText },
]

export function MorePage() {
  return (
    <DashboardLayout width="full">
      <section className="flex h-full min-h-0 w-full max-w-none flex-col overflow-hidden bg-[#eaf3fb]">
        <header className="grid w-full flex-none gap-4 px-4 pb-4 pt-[calc(16px+env(safe-area-inset-top))] sm:px-5 md:px-8 md:pt-6 lg:px-10">
          <h1 className="text-center text-[18px] font-semibold leading-tight text-[#18202f]">
            Mais Opcoes
          </h1>
        </header>

        <div className="scrollbar-none min-h-0 w-full flex-1 overflow-y-auto overflow-x-hidden rounded-t-[28px] bg-white pb-[var(--app-mobile-content-bottom)] pt-2 shadow-[0_-1px_8px_rgba(17,24,39,0.04)] md:pb-0">
          {links.map(({ to, label, icon: Icon }) => (
            <Link
              key={to}
              to={to}
              className="grid min-h-[54px] w-full cursor-pointer grid-cols-[32px_minmax(0,1fr)] items-center gap-3 border-b border-[#edf1f6] px-4 py-3 text-left transition-colors hover:bg-[#f8fafc] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#2563eb] md:px-8 lg:px-10"
            >
              <span className="grid h-8 w-8 place-items-center text-[#1f2937]">
                <Icon className="h-[21px] w-[21px]" aria-hidden="true" />
              </span>
              <span className="truncate text-[15px] font-medium text-[#1f2937]">
                {label}
              </span>
            </Link>
          ))}
        </div>
      </section>
    </DashboardLayout>
  )
}
