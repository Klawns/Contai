export type AuthPath =
  | '/'
  | '/login'
  | '/registro'
  | '/more'
  | '/reports'
  | '/accounts'
  | '/accounts/new'
  | '/accounts/edit'
  | '/transactions'
  | '/transactions/edit'
  | '/transactions/income/new'
  | '/transactions/expense/new'
  | '/transactions/transfer/new'
  | '/credit-cards'
  | '/credit-cards/new'
  | '/credit-cards/edit'
  | '/credit-cards/purchase'
  | '/credit-cards/invoices'
  | '/credit-cards/invoice'
  | '/credit-cards/invoice/pay'

export type NavigationPath = AuthPath | `${AuthPath}?${string}`

export const supportedPaths = new Set<string>([
  '/',
  '/login',
  '/registro',
  '/more',
  '/reports',
  '/accounts',
  '/accounts/new',
  '/accounts/edit',
  '/transactions',
  '/transactions/edit',
  '/transactions/income/new',
  '/transactions/expense/new',
  '/transactions/transfer/new',
  '/credit-cards',
  '/credit-cards/new',
  '/credit-cards/edit',
  '/credit-cards/purchase',
  '/credit-cards/invoices',
  '/credit-cards/invoice',
  '/credit-cards/invoice/pay',
])

export function getAuthPath(pathname: string): AuthPath {
  return supportedPaths.has(pathname) ? (pathname as AuthPath) : '/'
}
