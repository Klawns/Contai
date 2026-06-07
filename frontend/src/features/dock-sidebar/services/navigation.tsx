import {
  ArrowLeftRight,
  CalendarDays,
  Ellipsis,
  House,
  ReceiptText,
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
  {
    label: 'Planejamento',
    icon: CalendarDays,
    path: '/planning',
    children: [
      { label: 'Geral', path: '/planning' },
      { label: 'Conta a pagar', path: '/planning/payables/new' },
      { label: 'Conta a receber', path: '/planning/receivables/new' },
    ],
  },
  { label: 'Mais', icon: Ellipsis, path: '/more' },
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
      (item.path === '/planning' && currentPath.startsWith('/planning')),
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
    {
      label: 'Conta a pagar',
      icon: ReceiptText,
      color: '#d93658',
      onSelect: () => navigate('/planning/payables/new'),
    },
    {
      label: 'Conta a receber',
      icon: CalendarDays,
      color: '#159c57',
      onSelect: () => navigate('/planning/receivables/new'),
    },
  ]
}
