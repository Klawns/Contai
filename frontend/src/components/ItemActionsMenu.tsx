import { useEffect, useId, useRef, useState } from 'react'
import { Ellipsis } from 'lucide-react'

type ItemActionsMenuProps = {
  label: string
  onEdit: () => void
  onDelete: () => void
  isDeleteDisabled?: boolean
}

export function ItemActionsMenu({
  label,
  onEdit,
  onDelete,
  isDeleteDisabled = false,
}: ItemActionsMenuProps) {
  const [isOpen, setIsOpen] = useState(false)
  const menuId = useId()
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!isOpen) {
      return undefined
    }

    function handlePointerDown(event: PointerEvent) {
      if (!menuRef.current?.contains(event.target as Node)) {
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
    <div ref={menuRef} className="relative">
      <button
        type="button"
        className="grid h-8 w-8 cursor-pointer place-items-center rounded-full text-[#958c9f] transition-colors hover:bg-[#f4f0f8] hover:text-[#4d168f] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        aria-label={label}
        aria-haspopup="menu"
        aria-expanded={isOpen}
        aria-controls={menuId}
        onClick={() => setIsOpen((current) => !current)}
      >
        <Ellipsis className="h-5 w-5" aria-hidden="true" />
      </button>
      {isOpen ? (
        <div
          id={menuId}
          role="menu"
          className="absolute right-0 top-9 z-20 w-[150px] overflow-hidden rounded-xl border border-[#ece7f3] bg-white py-1.5 text-[#2c2237] shadow-[0_18px_45px_rgba(35,24,52,0.16)]"
        >
          <button
            type="button"
            role="menuitem"
            className="h-10 w-full cursor-pointer px-3 text-left text-[13px] font-semibold text-[#4f435c] transition-colors hover:bg-[#f8f5fb]"
            onClick={() => {
              setIsOpen(false)
              onEdit()
            }}
          >
            Editar
          </button>
          <button
            type="button"
            role="menuitem"
            className="h-10 w-full cursor-pointer px-3 text-left text-[13px] font-semibold text-[#c75959] transition-colors hover:bg-[#fff4f4] disabled:cursor-not-allowed disabled:opacity-55"
            disabled={isDeleteDisabled}
            onClick={() => {
              setIsOpen(false)
              onDelete()
            }}
          >
            Deletar
          </button>
        </div>
      ) : null}
    </div>
  )
}
