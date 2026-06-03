export type MonthlyDashboard = {
  userId: string
  period: DashboardPeriod
  totalBalance: number
  monthlyIncome: number
  monthlyExpense: number
  monthlyTransferIn: number
  monthlyTransferOut: number
  monthlyNetBalance: number
  accountBalances: AccountBalance[]
  expensesByCategory: ExpenseByCategory[]
  recentTransactions: RecentTransaction[]
}

export type DashboardPeriod = {
  startAt: string
  endAt: string
}

export type AccountBalance = {
  accountId: string
  name: string
  type: string
  balance: number
  bankIconId: string
}

export type ExpenseByCategory = {
  categoryId: string
  name: string
  color: string
  icon: string
  total: number
}

export type RecentTransaction = {
  id: string
  userId: string
  type: string
  description: string
  amount: number
  occurredAt: string
  accountId: string | null
  sourceAccountId: string | null
  destinationAccountId: string | null
  categoryId: string | null
  status: string
  note: string
  removedAt: string | null
  createdAt: string
  updatedAt: string
}

export type MonthlyFinancialSeriesPoint = {
  monthLabel: string
  income: number
  expense: number
  balance: number
}

export type MonthlySeries = {
  userId: string
  period: DashboardPeriod
  points: MonthlySeriesPoint[]
}

export type MonthlySeriesPoint = {
  monthStartAt: string
  monthEndAt: string
  income: number
  expense: number
  balance: number
}
