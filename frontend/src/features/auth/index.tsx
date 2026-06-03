import { useEffect, useState } from 'react'
import { Dashboard } from '../dashboard'
import { DockSidebarLayout } from '../dock-sidebar'
import { MorePage } from '../more'
import { LoginForm } from './components/LoginForm.tsx'
import { RegisterForm } from './components/RegisterForm.tsx'
import { AuthShell } from './components/AuthShell.tsx'
import { useLogout } from './hooks/useLogout.ts'
import { useCurrentUser } from './hooks/useCurrentUser.ts'
import {
  getAuthPath,
  navigateTo,
  navigationEventName,
  type AuthPath,
} from './services/navigation.ts'
import type { AuthenticatedUser } from './types/auth.ts'

function useAuthPath() {
  const [path, setPath] = useState(() => getAuthPath())

  useEffect(() => {
    const handleNavigation = () => setPath(getAuthPath())

    window.addEventListener('popstate', handleNavigation)
    window.addEventListener(navigationEventName, handleNavigation)

    return () => {
      window.removeEventListener('popstate', handleNavigation)
      window.removeEventListener(navigationEventName, handleNavigation)
    }
  }, [])

  return path
}

type AuthenticatedAppProps = {
  path: AuthPath
  user: AuthenticatedUser
}

function AuthenticatedApp({ path, user }: AuthenticatedAppProps) {
  const logoutMutation = useLogout()
  const handleLogout = () => logoutMutation.mutate()
  const content =
    path === '/more' ? (
      <MorePage isLoggingOut={logoutMutation.isPending} onLogout={handleLogout} />
    ) : (
      <Dashboard
        user={user}
        isLoggingOut={logoutMutation.isPending}
        onLogout={handleLogout}
      />
    )

  return (
    <DockSidebarLayout
      currentPath={path}
      isLoggingOut={logoutMutation.isPending}
      onLogout={handleLogout}
    >
      <div key={path} className="contents">
        {content}
      </div>
    </DockSidebarLayout>
  )
}

export function AuthGate() {
  const path = useAuthPath()
  const currentUserQuery = useCurrentUser()

  useEffect(() => {
    if (currentUserQuery.isLoading) {
      return
    }

    const hasUser = Boolean(currentUserQuery.data)
    const isAuthRoute = path === '/login' || path === '/registro'

    if (!isAuthRoute && !hasUser) {
      navigateTo('/login', { replace: true })
      return
    }

    if (isAuthRoute && hasUser) {
      navigateTo('/', { replace: true })
    }
  }, [currentUserQuery.data, currentUserQuery.isLoading, path])

  if (currentUserQuery.isLoading) {
    return (
      <main className="grid min-h-svh place-items-center bg-[#f4f7fb] px-4 text-[14px] font-medium text-[#6f6679]">
        Validando sessao...
      </main>
    )
  }

  if (path === '/login') {
    return (
      <AuthShell title="Entrar" subtitle="Acesse sua conta para acompanhar suas financas.">
        <LoginForm />
      </AuthShell>
    )
  }

  if (path === '/registro') {
    return (
      <AuthShell title="Criar conta" subtitle="Cadastre-se para organizar contas e saldos.">
        <RegisterForm />
      </AuthShell>
    )
  }

  if (!currentUserQuery.data) {
    return null
  }

  return <AuthenticatedApp path={path} user={currentUserQuery.data} />
}
