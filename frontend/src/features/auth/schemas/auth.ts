import { z } from 'zod'

export const loginSchema = z.object({
  email: z.string().trim().email('Informe um e-mail valido.'),
  password: z.string().min(1, 'Informe sua senha.'),
})

export const registerSchema = z.object({
  name: z.string().trim().min(2, 'Informe seu nome.'),
  email: z.string().trim().email('Informe um e-mail valido.'),
  password: z.string().min(8, 'A senha deve ter pelo menos 8 caracteres.'),
})

export const authenticatedUserSchema = z.object({
  id: z.string(),
  name: z.string().optional(),
  email: z.string().email(),
  status: z.string(),
})

export const createdUserSchema = authenticatedUserSchema.extend({
  name: z.string(),
  createdAt: z.string(),
})
