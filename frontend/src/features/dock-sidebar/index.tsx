import type { ReactNode } from 'react'
import type { AuthPath } from '../auth/services/navigation'
import { BottomDock } from './components/BottomDock'
import { QuickActionOverlay } from './components/QuickActionOverlay'
import { Sidebar } from './components/Sidebar'
import { useQuickActions } from './hooks/useQuickActions'
import { getNavigationItems, quickActions } from './services/navigation'

type DockSidebarLayoutProps = {
  children: ReactNode
  currentPath: AuthPath
  isLoggingOut: boolean
  onLogout: () => void
}

export function DockSidebarLayout({
  children,
  currentPath,
  isLoggingOut,
  onLogout,
}: DockSidebarLayoutProps) {
  const quickActionsState = useQuickActions()
  const navigationItems = getNavigationItems(currentPath)

  return (
    <div className="min-h-svh bg-[#f4f7fb] md:grid md:grid-cols-[232px_minmax(0,1fr)]">
      <Sidebar items={navigationItems} isLoggingOut={isLoggingOut} onLogout={onLogout} />
      <div className="relative min-h-svh md:min-w-0">
        {children}
        <QuickActionOverlay
          actions={quickActions}
          isOpen={quickActionsState.isOpen}
          onClose={quickActionsState.close}
        />
        <BottomDock
          items={navigationItems}
          isQuickActionsOpen={quickActionsState.isOpen}
          onToggleQuickActions={quickActionsState.toggle}
        />
      </div>
    </div>
  )
}
