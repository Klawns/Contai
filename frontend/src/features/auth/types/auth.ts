import type { z } from 'zod'
import type {
  authenticatedUserSchema,
  createdUserSchema,
  loginSchema,
  registerSchema,
} from '../schemas/auth.ts'

export type LoginPayload = z.infer<typeof loginSchema>
export type RegisterPayload = z.infer<typeof registerSchema>
export type AuthenticatedUser = z.infer<typeof authenticatedUserSchema>
export type CreatedUser = z.infer<typeof createdUserSchema>
