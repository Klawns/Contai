import {
  ArrowLeftRight,
  CalendarDays,
  Ellipsis,
  House,
  TrendingDown,
  TrendingUp,
  type LucideIcon,
} from 'lucide-react'
import { navigateTo, type AuthPath } from '../../auth/services/navigation'

type IconComponent = LucideIcon

export type NavigationItem = {
  label: string
  icon: IconComponent
  path?: AuthPath
  active?: boolean
  onSelect?: () => void
}

export type QuickAction = {
  label: string
  icon: IconComponent
  color: string
  onSelect?: () => void
}

const baseNavigationItems: NavigationItem[] = [
  { label: 'Principal', icon: House, path: '/' },
  { label: 'Transacoes', icon: ArrowLeftRight },
  { label: 'Planejamento', icon: CalendarDays },
  { label: 'Mais', icon: Ellipsis, path: '/more' },
]

export function getNavigationItems(currentPath: AuthPath): NavigationItem[] {
  return baseNavigationItems.map((item) => ({
    ...item,
    active: item.path === currentPath,
    onSelect: item.path ? () => navigateTo(item.path!) : undefined,
  }))
}

export const quickActions: QuickAction[] = [
  { label: 'Receita', icon: TrendingUp, color: '#2bbf6a' },
  { label: 'Transferencia', icon: ArrowLeftRight, color: '#4aa8e8' },
  { label: 'Despesa', icon: TrendingDown, color: '#ef4771' },
]
