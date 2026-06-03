import { LogoutActionButton } from '../auth/components/LogoutActionButton'
import { DashboardLayout } from '../dashboard/components'

type MorePageProps = {
  isLoggingOut: boolean
  onLogout: () => void
}

export function MorePage({ isLoggingOut, onLogout }: MorePageProps) {
  return (
    <DashboardLayout>
      <section className="flex min-h-[calc(100svh-170px)] flex-col gap-3">
        <h1 className="text-[22px] font-semibold leading-tight text-[#241932]">Mais</h1>
        <p className="text-[14px] font-medium text-[#6f6679]">Nada por enquanto.</p>
        <div className="mt-auto flex justify-end">
          <LogoutActionButton isLoggingOut={isLoggingOut} onLogout={onLogout} />
        </div>
      </section>
    </DashboardLayout>
  )
}
