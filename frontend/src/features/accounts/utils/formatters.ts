export const accountTypeLabels = {
  checking: 'Conta corrente',
  savings: 'Poupanca',
  digital: 'Conta digital',
  cash: 'Carteira',
  salary: 'Conta salario',
  investment: 'Investimentos',
  other: 'Outra',
} as const

const shortDateFormatter = new Intl.DateTimeFormat('pt-BR', {
  day: '2-digit',
  month: 'short',
})

export function formatShortDate(date: Date) {
  return shortDateFormatter.format(date).replace('.', '').toUpperCase()
}
