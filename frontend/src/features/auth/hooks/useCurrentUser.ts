import { useQuery } from '@tanstack/react-query'
import { getCurrentUser } from '../services/authService.ts'

export const authQueryKeys = {
  currentUser: ['auth', 'me'] as const,
}

export function useCurrentUser() {
  return useQuery({
    queryKey: authQueryKeys.currentUser,
    queryFn: getCurrentUser,
    retry: false,
  })
}
