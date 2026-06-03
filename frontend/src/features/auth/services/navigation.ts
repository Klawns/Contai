export type AuthPath = '/' | '/login' | '/registro' | '/more'

const supportedPaths = new Set<string>(['/', '/login', '/registro', '/more'])
export const navigationEventName = 'contai:navigation'

export function getAuthPath(pathname = window.location.pathname): AuthPath {
  return supportedPaths.has(pathname) ? (pathname as AuthPath) : '/'
}

export function navigateTo(path: AuthPath, options: { replace?: boolean } = {}) {
  const method = options.replace ? 'replaceState' : 'pushState'
  window.history[method](null, '', path)
  window.dispatchEvent(new Event(navigationEventName))
}
