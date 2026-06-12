export const reportsQueryKeys = {
  all: ['reports'] as const,
  accounts: () => [...reportsQueryKeys.all, 'accounts'] as const,
  categories: () => [...reportsQueryKeys.all, 'categories'] as const,
}
