import { useCallback, useEffect, useMemo, useRef, useState, type ReactNode } from 'react'
import {
  ConfirmDialogContext,
  type ConfirmDialogOptions,
  type PendingConfirmation,
} from './confirm-dialog-context'

export function ConfirmDialogProvider({ children }: { children: ReactNode }) {
  const [pendingConfirmation, setPendingConfirmation] = useState<PendingConfirmation | null>(null)
  const cancelButtonRef = useRef<HTMLButtonElement>(null)

  const close = useCallback(
    (result: boolean) => {
      pendingConfirmation?.resolve(result)
      setPendingConfirmation(null)
    },
    [pendingConfirmation],
  )

  const confirm = useCallback((options: ConfirmDialogOptions) => {
    return new Promise<boolean>((resolve) => {
      setPendingConfirmation({
        title: options.title,
        description: options.description,
        confirmLabel: options.confirmLabel ?? 'Confirmar',
        cancelLabel: options.cancelLabel ?? 'Cancelar',
        tone: options.tone ?? 'default',
        resolve,
      })
    })
  }, [])

  const value = useMemo(() => ({ confirm }), [confirm])

  useEffect(() => {
    if (!pendingConfirmation) {
      return undefined
    }

    const previousActiveElement = document.activeElement instanceof HTMLElement
      ? document.activeElement
      : null

    cancelButtonRef.current?.focus()

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        event.preventDefault()
        close(false)
      }
    }

    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      previousActiveElement?.focus()
    }
  }, [close, pendingConfirmation])

  return (
    <ConfirmDialogContext.Provider value={value}>
      {children}
      {pendingConfirmation ? (
        <div
          className="fixed inset-0 z-50 grid place-items-center bg-[#1f1828]/52 px-4 py-6 backdrop-blur-[2px]"
          role="presentation"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) {
              close(false)
            }
          }}
        >
          <section
            role="alertdialog"
            aria-modal="true"
            aria-labelledby="confirm-dialog-title"
            aria-describedby="confirm-dialog-description"
            className="w-full max-w-[380px] rounded-2xl border border-[#ece8f2] bg-white p-5 text-left shadow-[0_24px_70px_rgba(28,20,40,0.28)]"
          >
            <h2
              id="confirm-dialog-title"
              className="text-[18px] font-semibold leading-tight text-[#241a30]"
            >
              {pendingConfirmation.title}
            </h2>
            <p
              id="confirm-dialog-description"
              className="mt-2 text-[14px] font-medium leading-relaxed text-[#6f6679]"
            >
              {pendingConfirmation.description}
            </p>
            <div className="mt-5 grid gap-2 sm:flex sm:flex-row-reverse">
              <button
                type="button"
                className={`h-11 cursor-pointer rounded-full px-5 text-[14px] font-semibold text-white transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 ${
                  pendingConfirmation.tone === 'danger'
                    ? 'bg-[#c83b3b] hover:bg-[#ad2f2f] focus-visible:outline-[#c83b3b]'
                    : 'bg-[#6818e8] hover:bg-[#5712c9] focus-visible:outline-[#7b2cff]'
                }`}
                onClick={() => close(true)}
              >
                {pendingConfirmation.confirmLabel}
              </button>
              <button
                ref={cancelButtonRef}
                type="button"
                className="h-11 cursor-pointer rounded-full border border-[#e3ddea] bg-white px-5 text-[14px] font-semibold text-[#4f435c] transition-colors hover:bg-[#f8f5fb] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                onClick={() => close(false)}
              >
                {pendingConfirmation.cancelLabel}
              </button>
            </div>
          </section>
        </div>
      ) : null}
    </ConfirmDialogContext.Provider>
  )
}
