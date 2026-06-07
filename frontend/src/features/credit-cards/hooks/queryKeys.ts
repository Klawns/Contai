export const creditCardQueryKeys = {
  all: ['credit-cards'] as const,
  lists: () => [...creditCardQueryKeys.all, 'list'] as const,
  purchases: (cardId: string) => [...creditCardQueryKeys.all, 'purchases', cardId] as const,
  invoices: (cardId: string) => [...creditCardQueryKeys.all, 'invoices', cardId] as const,
  invoice: (invoiceId: string) => [...creditCardQueryKeys.all, 'invoice', invoiceId] as const,
}
