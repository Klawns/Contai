import type { MonthlyFinancialSeriesPoint } from '../types/dashboard.ts'

export const summaryDisplayState = {
  isBalanceHidden: false,
}

export const monthlyFinancialSeries: MonthlyFinancialSeriesPoint[] = [
  {
    monthLabel: 'Nov',
    income: 780000,
    expense: 1185000,
    balance: -405000,
  },
  {
    monthLabel: 'Dez',
    income: 910000,
    expense: 1264000,
    balance: -354000,
  },
  {
    monthLabel: 'Jan',
    income: 835000,
    expense: 1198000,
    balance: -363000,
  },
  {
    monthLabel: 'Fev',
    income: 872000,
    expense: 1280500,
    balance: -408500,
  },
  {
    monthLabel: 'Mar',
    income: 806000,
    expense: 1219000,
    balance: -413000,
  },
  {
    monthLabel: 'Abr',
    income: 842000,
    expense: 1323100,
    balance: -481100,
  },
]
