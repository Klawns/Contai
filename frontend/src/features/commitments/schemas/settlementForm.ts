import { z } from 'zod'

export const settlementFormSchema = z.object({
  amount: z.number().int().positive('Informe um valor maior que zero.'),
  occurredOn: z.string().min(1, 'Informe a data.'),
  accountId: z.string().min(1, 'Selecione uma conta.'),
  categoryId: z.string().min(1, 'Selecione uma categoria.'),
  note: z.string(),
})
