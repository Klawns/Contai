import { lazy, Suspense } from 'react'
import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { DockSidebarLayout } from '../dock-sidebar'
import { LoginForm } from './components/LoginForm.tsx'
import { RegisterForm } from './components/RegisterForm.tsx'
import { AuthShell } from './components/AuthShell.tsx'
import { useLogout } from './hooks/useLogout.ts'
import { useCurrentUser } from './hooks/useCurrentUser.ts'
import { getAuthPath } from './services/navigation.ts'

const Dashboard = lazy(() =>
  import('../dashboard').then((module) => ({ default: module.Dashboard })),
)
const MorePage = lazy(() =>
  import('../more').then((module) => ({ default: module.MorePage })),
)
const AccountListPage = lazy(() =>
  import('../accounts').then((module) => ({ default: module.AccountListPage })),
)
const AccountCreatePage = lazy(() =>
  import('../accounts').then((module) => ({ default: module.AccountCreatePage })),
)
const AccountEditPage = lazy(() =>
  import('../accounts').then((module) => ({ default: module.AccountEditPage })),
)
const TransactionListPage = lazy(() =>
  import('../transactions').then((module) => ({ default: module.TransactionListPage })),
)
const TransactionCreatePage = lazy(() =>
  import('../transactions').then((module) => ({ default: module.TransactionCreatePage })),
)
const TransactionEditPage = lazy(() =>
  import('../transactions').then((module) => ({ default: module.TransactionEditPage })),
)
const CreditCardListPage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.CreditCardListPage })),
)
const CreditCardFormPage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.CreditCardFormPage })),
)
const CardPurchasePage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.CardPurchasePage })),
)
const CardInvoicesPage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.CardInvoicesPage })),
)
const CardInvoiceDetailPage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.CardInvoiceDetailPage })),
)
const PayCardInvoicePage = lazy(() =>
  import('../credit-cards').then((module) => ({ default: module.PayCardInvoicePage })),
)

function AuthLoading() {
  return (
    <main className="grid min-h-svh place-items-center bg-[#f4f7fb] px-4 text-[14px] font-medium text-[#6f6679]">
      Validando sessao...
    </main>
  )
}

function ProtectedRoute() {
  const currentUserQuery = useCurrentUser()

  if (currentUserQuery.isLoading) {
    return <AuthLoading />
  }

  if (!currentUserQuery.data) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}

function PublicOnlyRoute() {
  const currentUserQuery = useCurrentUser()

  if (currentUserQuery.isLoading) {
    return <AuthLoading />
  }

  if (currentUserQuery.data) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}

function AuthenticatedLayout() {
  const location = useLocation()
  const currentUserQuery = useCurrentUser()
  const logoutMutation = useLogout()
  const currentPath = getAuthPath(location.pathname)
  const user = currentUserQuery.data

  if (!user) {
    return null
  }

  return (
    <DockSidebarLayout
      currentPath={currentPath}
      isLoggingOut={logoutMutation.isPending}
      onLogout={() => logoutMutation.mutate()}
    >
      <div key={`${location.pathname}${location.search}`} className="contents">
        <Outlet context={{ user, logoutMutation }} />
      </div>
    </DockSidebarLayout>
  )
}

function DashboardRoute() {
  const currentUserQuery = useCurrentUser()
  const logoutMutation = useLogout()
  const user = currentUserQuery.data

  if (!user) {
    return null
  }

  return (
    <Dashboard
      user={user}
      isLoggingOut={logoutMutation.isPending}
      onLogout={() => logoutMutation.mutate()}
    />
  )
}

function FallbackRoute() {
  const currentUserQuery = useCurrentUser()

  if (currentUserQuery.isLoading) {
    return <AuthLoading />
  }

  return <Navigate to={currentUserQuery.data ? '/' : '/login'} replace />
}

export function AuthGate() {
  return (
    <Suspense fallback={<AuthLoading />}>
      <Routes>
        <Route element={<PublicOnlyRoute />}>
          <Route
            path="/login"
            element={(
              <AuthShell title="Entrar" subtitle="Acesse sua conta para acompanhar suas financas.">
                <LoginForm />
              </AuthShell>
            )}
          />
          <Route
            path="/registro"
            element={(
              <AuthShell title="Criar conta" subtitle="Cadastre-se para organizar contas e saldos.">
                <RegisterForm />
              </AuthShell>
            )}
          />
        </Route>

        <Route element={<ProtectedRoute />}>
          <Route element={<AuthenticatedLayout />}>
            <Route index element={<DashboardRoute />} />
            <Route path="/more" element={<MorePage />} />
            <Route path="/accounts" element={<AccountListPage />} />
            <Route path="/accounts/new" element={<AccountCreatePage />} />
            <Route path="/accounts/edit" element={<AccountEditPage />} />
            <Route path="/transactions" element={<TransactionListPage />} />
            <Route path="/transactions/edit" element={<TransactionEditPage />} />
            <Route path="/transactions/income/new" element={<TransactionCreatePage type="income" />} />
            <Route path="/transactions/expense/new" element={<TransactionCreatePage type="expense" />} />
            <Route path="/transactions/transfer/new" element={<TransactionCreatePage type="transfer" />} />
            <Route path="/credit-cards" element={<CreditCardListPage />} />
            <Route path="/credit-cards/new" element={<CreditCardFormPage />} />
            <Route path="/credit-cards/edit" element={<CreditCardFormPage mode="edit" />} />
            <Route path="/credit-cards/purchase" element={<CardPurchasePage />} />
            <Route path="/credit-cards/invoices" element={<CardInvoicesPage />} />
            <Route path="/credit-cards/invoice" element={<CardInvoiceDetailPage />} />
            <Route path="/credit-cards/invoice/pay" element={<PayCardInvoicePage />} />
          </Route>
        </Route>

        <Route path="*" element={<FallbackRoute />} />
      </Routes>
    </Suspense>
  )
}
