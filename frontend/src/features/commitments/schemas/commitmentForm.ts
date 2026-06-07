import { z } from 'zod'

export const commitmentFormSchema = z
  .object({
    description: z.string().trim().min(1, 'Informe a descricao.'),
    amount: z.number().int().positive('Informe um valor maior que zero.'),
    dueOn: z.string().min(1, 'Informe o vencimento.'),
    accountId: z.string().min(1, 'Selecione uma conta.'),
    categoryId: z.string().min(1, 'Selecione uma categoria.'),
    note: z.string(),
    hasRecurrence: z.boolean(),
    recurrenceFrequency: z.enum(['daily', 'weekly', 'monthly', 'yearly']),
    recurrenceInterval: z.number().int().positive('Informe um intervalo maior que zero.'),
    recurrenceEndsOn: z.string(),
  })
  .superRefine((values, context) => {
    if (values.hasRecurrence && values.recurrenceEndsOn) {
      const dueDate = Date.parse(values.dueOn)
      const endDate = Date.parse(values.recurrenceEndsOn)

      if (Number.isFinite(dueDate) && Number.isFinite(endDate) && endDate < dueDate) {
        context.addIssue({
          code: 'custom',
          path: ['recurrenceEndsOn'],
          message: 'O fim da recorrencia deve ser apos o vencimento.',
        })
      }
    }
  })
