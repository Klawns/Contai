import type { ReactNode } from 'react'
import { useCallback, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import type { AuthPath } from '../auth/services/navigation'
import { BottomDock } from './components/BottomDock'
import { QuickActionOverlay } from './components/QuickActionOverlay'
import { Sidebar } from './components/Sidebar'
import { useQuickActions } from './hooks/useQuickActions'
import { getNavigationItems, getQuickActions } from './services/navigation'

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
  const navigate = useNavigate()
  const quickActionsState = useQuickActions()
  const goToPath = useCallback((path: AuthPath) => navigate(path), [navigate])
  const navigationItems = useMemo(
    () => getNavigationItems(currentPath, goToPath),
    [currentPath, goToPath],
  )
  const quickActions = useMemo(() => getQuickActions(goToPath), [goToPath])

  return (
    <div className="h-[var(--app-viewport-height)] min-h-0 overflow-hidden bg-[#f4f7fb] md:grid md:grid-cols-[232px_minmax(0,1fr)]">
      <Sidebar items={navigationItems} isLoggingOut={isLoggingOut} onLogout={onLogout} />
      <div className="relative h-full min-h-0 overflow-hidden md:min-w-0">
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
