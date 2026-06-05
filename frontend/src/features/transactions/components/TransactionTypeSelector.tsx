import { useEffect, useId, useRef, useState } from 'react'
import { Check, ChevronDown } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import type { TransactionType } from '../types/transactions.ts'

type TransactionTypeSelectorProps = {
  type: TransactionType
  isLocked?: boolean
}

const options = [
  {
    type: 'income',
    label: 'Receita',
    path: '/transactions/income/new',
  },
  {
    type: 'expense',
    label: 'Despesa',
    path: '/transactions/expense/new',
  },
  {
    type: 'transfer',
    label: 'Transferencia',
    path: '/transactions/transfer/new',
  },
] as const

export function TransactionTypeSelector({ type, isLocked = false }: TransactionTypeSelectorProps) {
  const navigate = useNavigate()
  const [isOpen, setIsOpen] = useState(false)
  const menuId = useId()
  const selectorRef = useRef<HTMLDivElement>(null)
  const selected = options.find((option) => option.type === type) ?? options[1]

  useEffect(() => {
    if (!isOpen) {
      return undefined
    }

    function handlePointerDown(event: PointerEvent) {
      if (!selectorRef.current?.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('pointerdown', handlePointerDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen])

  return (
    <div ref={selectorRef} className="relative z-20 flex justify-center">
      <button
        type="button"
        className={`inline-flex h-10 items-center gap-2 rounded-full bg-black/18 px-4 text-[14px] font-semibold text-white shadow-[0_10px_24px_rgba(0,0,0,0.1)] backdrop-blur transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white ${
          isLocked ? 'cursor-default' : 'cursor-pointer hover:bg-black/24'
        }`}
        aria-haspopup="menu"
        aria-expanded={!isLocked && isOpen}
        aria-controls={menuId}
        onClick={() => {
          if (!isLocked) {
            setIsOpen((current) => !current)
          }
        }}
      >
        {selected.label}
        {isLocked ? null : (
          <ChevronDown
            className={`h-4 w-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
            aria-hidden="true"
          />
        )}
      </button>

      {!isLocked && isOpen ? (
        <div
          id={menuId}
          role="menu"
          className="absolute top-12 w-[190px] overflow-hidden rounded-2xl bg-white py-1.5 text-[#2c2237] shadow-[0_18px_45px_rgba(35,24,52,0.24)] ring-1 ring-black/6"
        >
          {options.map((option) => (
            <button
              key={option.type}
              type="button"
              role="menuitem"
              className={`flex h-11 w-full cursor-pointer items-center justify-between px-4 text-left text-[14px] font-semibold transition-colors ${
                option.type === type
                  ? 'bg-[#f4f1f7] text-[#281d35]'
                  : 'text-[#5f536d] hover:bg-[#f8f5fb]'
              }`}
              onClick={() => {
                setIsOpen(false)
                if (option.type !== type) {
                  navigate(option.path)
                }
              }}
            >
              {option.label}
              {option.type === type ? <Check className="h-4 w-4" aria-hidden="true" /> : null}
            </button>
          ))}
        </div>
      ) : null}
    </div>
  )
}
