import {
  ArrowLeftRight,
  CreditCard,
  Ellipsis,
  House,
  TrendingDown,
  TrendingUp,
  type LucideIcon,
} from 'lucide-react'
import type { AuthPath } from '../../auth/services/navigation'

type IconComponent = LucideIcon

export type NavigationItem = {
  label: string
  icon: IconComponent
  path?: AuthPath
  active?: boolean
  onSelect?: () => void
  children?: NavigationSubItem[]
}

export type NavigationSubItem = {
  label: string
  path: AuthPath
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
  {
    label: 'Transacoes',
    icon: ArrowLeftRight,
    path: '/transactions',
    children: [
      { label: 'Geral', path: '/transactions' },
      { label: 'Receita', path: '/transactions/income/new' },
      { label: 'Despesa', path: '/transactions/expense/new' },
      { label: 'Transferencia', path: '/transactions/transfer/new' },
    ],
  },
  { label: 'Mais', icon: Ellipsis, path: '/more' },
  { label: 'Cartoes', icon: CreditCard, path: '/credit-cards' },
]

export function getNavigationItems(
  currentPath: AuthPath,
  navigate: (path: AuthPath) => void,
): NavigationItem[] {
  return baseNavigationItems.map((item) => ({
    ...item,
    active:
      item.path === currentPath ||
      (item.path === '/transactions' && currentPath.startsWith('/transactions')) ||
      (item.path === '/credit-cards' && currentPath.startsWith('/credit-cards')),
    onSelect: item.path ? () => navigate(item.path!) : undefined,
    children: item.children?.map((child) => ({
      ...child,
      active: child.path === currentPath,
      onSelect: () => navigate(child.path),
    })),
  }))
}

export function getQuickActions(navigate: (path: AuthPath) => void): QuickAction[] {
  return [
    {
      label: 'Receita',
      icon: TrendingUp,
      color: '#2bbf6a',
      onSelect: () => navigate('/transactions/income/new'),
    },
    {
      label: 'Transferencia',
      icon: ArrowLeftRight,
      color: '#4aa8e8',
      onSelect: () => navigate('/transactions/transfer/new'),
    },
    {
      label: 'Despesa',
      icon: TrendingDown,
      color: '#ef4771',
      onSelect: () => navigate('/transactions/expense/new'),
    },
  ]
}
