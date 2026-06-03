export function getBalanceTone(valueInCents: number) {
  if (valueInCents < 0) {
    return {
      state: 'negative',
      textClass: 'text-[#d83b3b]',
    } as const
  }

  if (valueInCents === 0) {
    return {
      state: 'zero',
      textClass: 'text-[#241a30]',
    } as const
  }

  return {
    state: 'positive',
    textClass: 'text-[#159b58]',
  } as const
}
