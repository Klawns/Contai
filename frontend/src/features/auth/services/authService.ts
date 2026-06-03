import { api } from '../../../lib/api/axios.ts'
import {
  authenticatedUserSchema,
  createdUserSchema,
  loginSchema,
  registerSchema,
} from '../schemas/auth.ts'
import type { AuthenticatedUser, CreatedUser, LoginPayload, RegisterPayload } from '../types/auth.ts'

export async function login(payload: LoginPayload): Promise<AuthenticatedUser> {
  const body = loginSchema.parse(payload)
  const response = await api.post<unknown>('/auth/login', body)

  return authenticatedUserSchema.parse(response.data)
}

export async function register(payload: RegisterPayload): Promise<CreatedUser> {
  const body = registerSchema.parse(payload)
  const response = await api.post<unknown>('/users', body)

  return createdUserSchema.parse(response.data)
}

export async function logout(): Promise<void> {
  await api.post('/auth/logout')
}

export async function getCurrentUser(): Promise<AuthenticatedUser> {
  const response = await api.get<unknown>('/auth/me')

  return authenticatedUserSchema.parse(response.data)
}
