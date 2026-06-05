import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { dashboardQueryKeys } from '../../dashboard/hooks/useMonthlyDashboard.ts'
import { logout } from '../services/authService.ts'
import { authQueryKeys } from './useCurrentUser.ts'

export function useLogout() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  return useMutation({
    mutationFn: logout,
    onSettled: () => {
      queryClient.removeQueries({ queryKey: dashboardQueryKeys.all })
      queryClient.setQueryData(authQueryKeys.currentUser, null)
      navigate('/login', { replace: true })
    },
  })
}
