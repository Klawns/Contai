export type AuthPath =
  | '/'
  | '/login'
  | '/registro'
  | '/more'
  | '/accounts'
  | '/accounts/new'
  | '/accounts/edit'
  | '/transactions'
  | '/transactions/edit'
  | '/transactions/income/new'
  | '/transactions/expense/new'
  | '/transactions/transfer/new'

export type NavigationPath = AuthPath | `${AuthPath}?${string}`

export const supportedPaths = new Set<string>([
  '/',
  '/login',
  '/registro',
  '/more',
  '/accounts',
  '/accounts/new',
  '/accounts/edit',
  '/transactions',
  '/transactions/edit',
  '/transactions/income/new',
  '/transactions/expense/new',
  '/transactions/transfer/new',
])

export function getAuthPath(pathname: string): AuthPath {
  return supportedPaths.has(pathname) ? (pathname as AuthPath) : '/'
}
