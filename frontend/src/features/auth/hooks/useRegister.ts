import { useMutation, useQueryClient } from '@tanstack/react-query'
import { authQueryKeys } from './useCurrentUser.ts'
import { register } from '../services/authService.ts'
import type { AuthenticatedUser, RegisterPayload } from '../types/auth.ts'

export function useRegister() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: RegisterPayload) => register(payload),
    onSuccess: (user) => {
      const authenticatedUser: AuthenticatedUser = {
        id: user.id,
        email: user.email,
        status: user.status,
      }

      queryClient.setQueryData(authQueryKeys.currentUser, authenticatedUser)
    },
  })
}
