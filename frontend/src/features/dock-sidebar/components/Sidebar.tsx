import { UserRound } from 'lucide-react'
import { LogoutActionButton } from '../../auth/components/LogoutActionButton'
import type { NavigationItem } from '../services/navigation'
import { DockItem } from './DockItem'

type SidebarProps = {
  items: NavigationItem[]
  isLoggingOut: boolean
  onLogout: () => void
}

export function Sidebar({ items, isLoggingOut, onLogout }: SidebarProps) {
  return (
    <aside
      className="sticky top-0 hidden h-svh flex-col gap-6 border-r border-[#e6e0ee] bg-white px-4 py-6 md:flex"
      aria-label="Navegacao principal"
    >
      <div className="flex items-center gap-2.5 text-lg text-[#241932]">
        <span className="grid h-[34px] w-[34px] place-items-center rounded-full border border-[#e4dfec] bg-[#fbfafe] text-[#6b6178]">
          <UserRound className="h-5 w-5" aria-hidden="true" />
        </span>
        <strong className="font-medium">Contai</strong>
      </div>
      <nav className="grid gap-1.5">
        {items.map((item) => (
          <DockItem key={item.label} item={item} variant="sidebar" />
        ))}
      </nav>
      <div className="mt-auto">
        <LogoutActionButton
          variant="sidebar"
          isLoggingOut={isLoggingOut}
          onLogout={onLogout}
        />
      </div>
    </aside>
  )
}
