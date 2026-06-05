import { useMemo } from 'react'
import { Landmark } from 'lucide-react'
import { svgBanco } from '@edusites/bancos-brasil'
import { supportedBankIds } from '../data/banks.ts'

type BankIconProps = {
  bankIconId?: string | null
  size?: number
  className?: string
}

export function BankIcon({ bankIconId, size = 40, className = '' }: BankIconProps) {
  const svg = useMemo(() => {
    if (!bankIconId || !supportedBankIds.has(bankIconId)) {
      return null
    }

    return svgBanco({ nome: bankIconId, formato: 'circulo', tamanho: size })
  }, [bankIconId, size])

  if (!svg) {
    return (
      <span
        className={`grid place-items-center rounded-full bg-[#f2eff8] text-[#6a22e5] ${className}`}
        style={{ width: size, height: size }}
      >
        <Landmark className="h-5 w-5" aria-hidden="true" />
      </span>
    )
  }

  return (
    <span
      className={`block overflow-hidden rounded-full ${className}`}
      style={{ width: size, height: size }}
      aria-hidden="true"
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  )
}
