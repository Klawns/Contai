import { useMutation, useQueryClient } from '@tanstack/react-query'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { logout } from '../services/authService.ts'
import { navigateTo } from '../services/navigation.ts'
import { authQueryKeys } from './useCurrentUser.ts'

export function useLogout() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: logout,
    onSettled: () => {
      queryClient.removeQueries({ queryKey: dashboardQueryKeys.all })
      queryClient.setQueryData(authQueryKeys.currentUser, null)
      navigateTo('/login', { replace: true })
    },
  })
}
