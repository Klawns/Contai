import { useState } from 'react'

export function useQuickActions() {
  const [isOpen, setIsOpen] = useState(false)

  return {
    isOpen,
    close: () => setIsOpen(false),
    toggle: () => setIsOpen((currentIsOpen) => !currentIsOpen),
  }
}
