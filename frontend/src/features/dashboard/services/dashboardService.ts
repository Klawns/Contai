import { api } from '../../../lib/api/axios.ts'
import {
  monthlyDashboardFiltersSchema,
  monthlyDashboardSchema,
  monthlySeriesSchema,
} from '../schemas/dashboard.ts'
import type { DashboardPeriod, MonthlyDashboard, MonthlySeries } from '../types/dashboard.ts'

export async function getMonthlyDashboard(
  filters: DashboardPeriod,
): Promise<MonthlyDashboard> {
  const params = monthlyDashboardFiltersSchema.parse(filters)
  const response = await api.get<unknown>('/dashboard/monthly', { params })

  return monthlyDashboardSchema.parse(response.data)
}

export async function getMonthlySeries(filters: DashboardPeriod): Promise<MonthlySeries> {
  const params = monthlyDashboardFiltersSchema.parse(filters)
  const response = await api.get<unknown>('/dashboard/monthly-series', { params })

  return monthlySeriesSchema.parse(response.data)
}
