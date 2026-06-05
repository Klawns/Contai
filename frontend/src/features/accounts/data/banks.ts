export type SupportedBank = {
  id: string
  label: string
}

export const supportedBanks = [
  { id: 'nubank', label: 'Nubank' },
  { id: 'cora', label: 'Cora' },
  { id: 'itau', label: 'Itau' },
  { id: 'inter', label: 'Inter' },
  { id: 'bancodobrasil', label: 'Banco do Brasil' },
  { id: 'bradesco', label: 'Bradesco' },
  { id: 'santander', label: 'Santander' },
  { id: 'caixa', label: 'Caixa' },
  { id: 'btg', label: 'BTG' },
  { id: 'xp', label: 'XP' },
  { id: 'picpay', label: 'PicPay' },
  { id: 'mercadopago', label: 'Mercado Pago' },
  { id: 'pagbank', label: 'PagBank' },
  { id: 'c6', label: 'C6 Bank' },
  { id: 'neon', label: 'Neon' },
  { id: 'sicoob', label: 'Sicoob' },
  { id: 'wise', label: 'Wise' },
  { id: 'paypal', label: 'PayPal' },
  { id: 'stone', label: 'Stone' },
  { id: 'next', label: 'Next' },
  { id: 'original', label: 'Original' },
  { id: 'sicredi', label: 'Sicredi' },
] as const satisfies readonly SupportedBank[]

export const supportedBankIds = new Set<string>(supportedBanks.map((bank) => bank.id))

export function getBankLabel(bankIconId: string) {
  return supportedBanks.find((bank) => bank.id === bankIconId)?.label ?? 'Banco'
}
