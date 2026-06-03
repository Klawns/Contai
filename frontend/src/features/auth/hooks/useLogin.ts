import { useMutation, useQueryClient } from '@tanstack/react-query'
import { authQueryKeys } from './useCurrentUser.ts'
import { login } from '../services/authService.ts'
import type { AuthenticatedUser, LoginPayload } from '../types/auth.ts'

export function useLogin() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: LoginPayload) => login(payload),
    onSuccess: (user: AuthenticatedUser) => {
      queryClient.setQueryData(authQueryKeys.currentUser, user)
    },
  })
}
