import { createContext, useContext } from 'react'

export type ConfirmDialogTone = 'default' | 'danger'

export type ConfirmDialogOptions = {
  title: string
  description: string
  confirmLabel?: string
  cancelLabel?: string
  tone?: ConfirmDialogTone
}

export type PendingConfirmation = Required<ConfirmDialogOptions> & {
  resolve: (value: boolean) => void
}

export type ConfirmDialogContextValue = {
  confirm: (options: ConfirmDialogOptions) => Promise<boolean>
}

export const ConfirmDialogContext = createContext<ConfirmDialogContextValue | null>(null)

export function useConfirmDialog() {
  const context = useContext(ConfirmDialogContext)

  if (!context) {
    throw new Error('useConfirmDialog must be used within ConfirmDialogProvider')
  }

  return context
}
