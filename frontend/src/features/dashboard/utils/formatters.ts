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

export function formatCurrencyOrHidden(valueInCents: number, isHidden = false) {
  if (isHidden) {
    return '---'
  }

  return formatCurrency(valueInCents)
}
