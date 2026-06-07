import { z } from 'zod'

const rfc3339Pattern =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/

export const rfc3339DateTimeSchema = z
  .string()
  .regex(rfc3339Pattern, 'Use uma data RFC3339 valida.')
  .refine((value) => !Number.isNaN(Date.parse(value)), {
    message: 'Use uma data RFC3339 valida.',
  })
