import { useQuery } from '@tanstack/react-query'
import { getMonthlyDashboard, getMonthlySeries } from '../services/dashboardService.ts'
import type { DashboardPeriod } from '../types/dashboard.ts'

export const dashboardQueryKeys = {
  all: ['dashboard'] as const,
  monthly: (period: DashboardPeriod) =>
    [...dashboardQueryKeys.all, 'monthly', period.startAt, period.endAt] as const,
  monthlySeries: (period: DashboardPeriod) =>
    [...dashboardQueryKeys.all, 'monthly-series', period.startAt, period.endAt] as const,
}

export function useMonthlyDashboard(period: DashboardPeriod) {
  return useQuery({
    queryKey: dashboardQueryKeys.monthly(period),
    queryFn: () => getMonthlyDashboard(period),
  })
}

export function useMonthlySeries(period: DashboardPeriod) {
  return useQuery({
    queryKey: dashboardQueryKeys.monthlySeries(period),
    queryFn: () => getMonthlySeries(period),
  })
}
