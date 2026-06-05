import { useMemo, useState } from 'react'
import { ChevronRight, Landmark } from 'lucide-react'
import { getBankLabel, supportedBanks } from '../data/banks.ts'
import { SelectionSheet } from '../../transactions/components/SelectionSheet.tsx'
import { BankIcon } from './BankIcon.tsx'

type BankIconSelectorProps = {
  value: string
  error?: string
  onChange: (bankIconId: string) => void
}

const quickBankIds = [
  'bradesco',
  'bancodobrasil',
  'santander',
  'itau',
  'nubank',
  'caixa',
  'inter',
]

export function BankIconSelector({ value, error, onChange }: BankIconSelectorProps) {
  const [isOpen, setIsOpen] = useState(false)
  const selectedBank = useMemo(
    () => supportedBanks.find((bank) => bank.id === value),
    [value],
  )
  const quickBanks = quickBankIds
    .map((bankIconId) => supportedBanks.find((bank) => bank.id === bankIconId))
    .filter((bank): bank is (typeof supportedBanks)[number] => Boolean(bank))

  return (
    <>
      <section className="grid gap-3 border-b border-[#f0ebf5] px-4 py-2.5 md:px-5">
        {selectedBank ? (
          <button
            type="button"
            className="relative min-h-12 w-full cursor-pointer text-left focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
            onClick={() => setIsOpen(true)}
          >
            <span className="absolute left-0 top-1/2 z-10 grid h-10 w-10 -translate-y-1/2 place-items-center">
              <BankIcon bankIconId={selectedBank.id} size={34} />
            </span>
            <span className="pointer-events-none absolute left-[52px] right-[52px] top-0 flex h-full items-center overflow-hidden">
              <span className="truncate text-left text-[16px] font-normal text-[#2c2237]">
                {selectedBank.label}
              </span>
            </span>
            <span className="absolute right-0 top-1/2 z-10 grid h-10 w-[52px] -translate-y-1/2 place-items-center text-[#9a91a5]">
              <ChevronRight className="h-5 w-5" aria-hidden="true" />
            </span>
          </button>
        ) : (
          <div className="grid gap-3">
            <h3 className="text-[14px] font-semibold text-[#2c2237]">
              Instituicao financeira
            </h3>

            <div className="grid w-full grid-cols-4 justify-items-center gap-y-4 md:hidden">
              {quickBanks.map((bank) => (
                <button
                  key={bank.id}
                  type="button"
                  className="grid h-16 w-16 cursor-pointer place-items-center rounded-full transition-transform hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7b2cff]"
                  aria-label={bank.label}
                  onClick={() => onChange(bank.id)}
                >
                  <BankIcon bankIconId={bank.id} size={64} />
                </button>
              ))}
              <button
                type="button"
                className="grid h-16 w-16 cursor-pointer place-items-center rounded-full bg-[#9f9da3] text-white transition-colors hover:bg-[#8f8d94] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7b2cff]"
                aria-label="Outros bancos"
                onClick={() => setIsOpen(true)}
              >
                <span className="text-[13px] font-semibold">Outros</span>
              </button>
            </div>

            <div className="hidden w-full grid-cols-[repeat(auto-fill,64px)] justify-start gap-x-5 gap-y-4 md:grid">
              {supportedBanks.map((bank) => (
                <button
                  key={bank.id}
                  type="button"
                  className="grid h-16 w-16 cursor-pointer place-items-center rounded-full transition-transform hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7b2cff]"
                  aria-label={bank.label}
                  onClick={() => onChange(bank.id)}
                >
                  <BankIcon bankIconId={bank.id} size={64} />
                </button>
              ))}
              <button
                type="button"
                className="grid h-16 w-16 cursor-pointer place-items-center rounded-full bg-[#9f9da3] text-white transition-colors hover:bg-[#8f8d94] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7b2cff]"
                aria-label="Outros bancos"
                onClick={() => setIsOpen(true)}
              >
                <span className="text-[13px] font-semibold">Outros</span>
              </button>
            </div>

            <button
              type="button"
              className="relative mt-1 min-h-12 w-full cursor-pointer text-left focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
              onClick={() => setIsOpen(true)}
            >
              <span className="absolute left-0 top-1/2 z-10 grid h-10 w-10 -translate-y-1/2 place-items-center text-[#6f647b]">
                <Landmark className="h-5 w-5" aria-hidden="true" />
              </span>
              <span className="pointer-events-none absolute left-[52px] right-[52px] top-0 flex h-full items-center overflow-hidden">
                <span className="truncate text-left text-[16px] font-normal text-[#aaa2b4]">
                  Selecione uma instituicao financeira
                </span>
              </span>
              <span className="absolute right-0 top-1/2 z-10 grid h-10 w-[52px] -translate-y-1/2 place-items-center text-[#9a91a5]">
                <ChevronRight className="h-5 w-5" aria-hidden="true" />
              </span>
            </button>
          </div>
        )}

        {error ? <p className="text-[12px] font-semibold text-[#c72f4d]">{error}</p> : null}
      </section>

      <SelectionSheet title="Instituicao financeira" isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="grid gap-2">
          {supportedBanks.map((bank) => {
            const isSelected = bank.id === value

            return (
              <button
                key={bank.id}
                type="button"
                className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
                  isSelected
                    ? 'border-[#7b2cff] bg-[#f7f2ff]'
                    : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
                }`}
                aria-pressed={isSelected}
                onClick={() => {
                  onChange(bank.id)
                  setIsOpen(false)
                }}
              >
                <BankIcon bankIconId={bank.id} size={40} />
                <span className="min-w-0 flex-1 truncate text-[14px] font-semibold text-[#2c2237]">
                  {getBankLabel(bank.id)}
                </span>
              </button>
            )
          })}
        </div>
      </SelectionSheet>
    </>
  )
}
