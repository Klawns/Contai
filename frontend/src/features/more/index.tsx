import { ChevronDown, CreditCard, FileText, Landmark } from 'lucide-react'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { DashboardLayout } from '../dashboard/components'
import { ReportDownloadPanel } from '../reports/components/ReportDownloadPanel.tsx'

export function MorePage() {
  const [isReportsOpen, setIsReportsOpen] = useState(false)
  const shouldReduceMotion = useReducedMotion()

  return (
    <DashboardLayout width="full">
      <section className="flex min-h-svh w-full max-w-none flex-col bg-[#eaf3fb] md:overflow-hidden">
        <header className="grid w-full gap-4 px-4 pb-4 pt-[calc(16px+env(safe-area-inset-top))] sm:px-5 md:px-8 md:pt-6 lg:px-10">
          <h1 className="text-center text-[18px] font-semibold leading-tight text-[#18202f]">
            Mais Opcoes
          </h1>

          <div className="grid w-full grid-cols-3 rounded-[12px] bg-[#dbe5ef] p-1 text-center text-[13px] font-semibold text-[#7c8795]">
            <span className="rounded-[9px] bg-white px-2 py-2 text-[#171d29] shadow-[0_1px_3px_rgba(23,29,41,0.08)]">
              Gerenciar
            </span>
            <span className="px-2 py-2">Acompanhar</span>
            <span className="px-2 py-2">Sobre</span>
          </div>
        </header>

        <div className="min-h-0 w-full flex-1 overflow-hidden rounded-t-[28px] bg-white pb-[calc(92px+env(safe-area-inset-bottom))] pt-2 shadow-[0_-1px_8px_rgba(17,24,39,0.04)] md:pb-0">
          <Link
            to="/accounts"
            className="grid min-h-[54px] w-full cursor-pointer grid-cols-[32px_minmax(0,1fr)] items-center gap-3 border-b border-[#edf1f6] px-4 py-3 text-left transition-colors hover:bg-[#f8fafc] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#2563eb] md:px-8 lg:px-10"
          >
            <span className="grid h-8 w-8 place-items-center text-[#1f2937]">
              <Landmark className="h-[21px] w-[21px]" aria-hidden="true" />
            </span>
            <span className="truncate text-[15px] font-medium text-[#1f2937]">Contas</span>
          </Link>

          <Link
            to="/credit-cards"
            className="grid min-h-[54px] w-full cursor-pointer grid-cols-[32px_minmax(0,1fr)] items-center gap-3 border-b border-[#edf1f6] px-4 py-3 text-left transition-colors hover:bg-[#f8fafc] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#2563eb] md:px-8 lg:px-10"
          >
            <span className="grid h-8 w-8 place-items-center text-[#1f2937]">
              <CreditCard className="h-[21px] w-[21px]" aria-hidden="true" />
            </span>
            <span className="truncate text-[15px] font-medium text-[#1f2937]">Cartoes</span>
          </Link>

          <button
            type="button"
            className="grid min-h-[54px] w-full cursor-pointer grid-cols-[32px_minmax(0,1fr)_20px] items-center gap-3 px-4 py-3 text-left transition-colors hover:bg-[#f8fafc] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#2563eb] md:px-8 lg:px-10"
            aria-expanded={isReportsOpen}
            onClick={() => setIsReportsOpen((current) => !current)}
          >
            <span className="grid h-8 w-8 place-items-center text-[#1f2937]">
              <FileText className="h-[21px] w-[21px]" aria-hidden="true" />
            </span>
            <span className="truncate text-[15px] font-medium text-[#1f2937]">Relatorios</span>
            <ChevronDown
              className={`h-5 w-5 text-[#9aa5b1] transition-transform ${
                isReportsOpen ? 'rotate-180' : ''
              }`}
              aria-hidden="true"
            />
          </button>

          <AnimatePresence initial={false}>
            {isReportsOpen ? (
              <motion.div
                className="overflow-hidden"
                initial={{
                  height: 0,
                  opacity: 0,
                  y: shouldReduceMotion ? 0 : -4,
                }}
                animate={{ height: 'auto', opacity: 1, y: 0 }}
                exit={{
                  height: 0,
                  opacity: 0,
                  y: shouldReduceMotion ? 0 : -4,
                }}
                transition={{ duration: shouldReduceMotion ? 0 : 0.16, ease: 'easeOut' }}
              >
                <ReportDownloadPanel />
              </motion.div>
            ) : null}
          </AnimatePresence>
        </div>
      </section>
    </DashboardLayout>
  )
}
