const currencyFormatter = new Intl.NumberFormat('pt-BR', {
  style: 'currency',
  currency: 'BRL',
  minimumFractionDigits: 2,
})

export function formatCurrency(valueInCents: number) {
  const absolute = Math.abs(valueInCents) / 100
  const formatted = currencyFormatter.format(absolute)

  return valueInCents < 0 ? `-${formatted}` : formatted
}

export function parseCurrencyToCents(value: string) {
  const digits = value.replace(/\D/g, '')

  return digits ? Number(digits) : 0
}

export function formatCentsAsInput(valueInCents: number) {
  return currencyFormatter.format(Math.max(0, valueInCents) / 100)
}
