import { useQuery } from '@tanstack/react-query'
import { listActiveCategories } from '../../transactions/services/categoryService.ts'
import { uniqueReportCategories } from '../utils/reportFilters.ts'
import { reportsQueryKeys } from './queryKeys.ts'

export function useReportCategories() {
  return useQuery({
    queryKey: reportsQueryKeys.categories(),
    queryFn: async () => {
      const [income, expense] = await Promise.all([
        listActiveCategories('income'),
        listActiveCategories('expense'),
      ])

      return [...income, ...expense]
    },
    select: uniqueReportCategories,
  })
}
