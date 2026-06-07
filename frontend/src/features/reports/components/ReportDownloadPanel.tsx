import { useState } from 'react'
import type { ReactNode } from 'react'
import { z } from 'zod'
import { CalendarDays, FileText, Landmark, ListFilter } from 'lucide-react'
import { MonthSelector, type SelectedMonth } from '../../../components/MonthSelector.tsx'
import { DateInput, TransactionFieldRow } from '../../transactions/components/FormFields.tsx'
import { AccountSelector } from '../../transactions/components/Selectors.tsx'
import { SelectionSheet } from '../../transactions/components/SelectionSheet.tsx'
import {
  formatLocalRFC3339,
  fromDateInputValue,
  getCurrentSelectedMonth,
  getMonthPeriod,
  toDateInputValue,
} from '../../transactions/utils/date.ts'
import type { ReportTransactionType } from '../types/reports.ts'
import {
  downloadAccountReportPDF,
  downloadAccountsReportPDF,
  downloadMonthlyReportPDF,
  downloadPeriodReportPDF,
  downloadTransactionsReportPDF,
  savePdfResponse,
} from '../services/reportService.ts'

type ReportID = 'accounts' | 'transactions' | 'period' | 'monthly' | 'account'

type ReportItem = {
  id: ReportID
  title: string
  description: string
  icon: ReactNode
}

type PeriodForm = {
  startOn: string
  endOn: string
}

type TransactionsForm = PeriodForm & {
  type: ReportTransactionType
}

type AccountForm = PeriodForm & {
  accountId: string
}

const today = toDateInputValue(new Date())

const reportItems: ReportItem[] = [
  {
    id: 'accounts',
    title: 'Relatorio de contas',
    description: 'Baixar PDF',
    icon: <FileText className="h-[18px] w-[18px]" aria-hidden="true" />,
  },
  {
    id: 'transactions',
    title: 'Receitas/despesas por periodo',
    description: 'Filtrar e baixar PDF',
    icon: <ListFilter className="h-[18px] w-[18px]" aria-hidden="true" />,
  },
  {
    id: 'period',
    title: 'Geral por periodo',
    description: 'Filtrar e baixar PDF',
    icon: <CalendarDays className="h-[18px] w-[18px]" aria-hidden="true" />,
  },
  {
    id: 'monthly',
    title: 'Mensal consolidado',
    description: 'Selecionar mes',
    icon: <FileText className="h-[18px] w-[18px]" aria-hidden="true" />,
  },
  {
    id: 'account',
    title: 'Por conta bancaria',
    description: 'Selecionar conta e periodo',
    icon: <Landmark className="h-[18px] w-[18px]" aria-hidden="true" />,
  },
]

const basePeriodSchema = z.object({
  startOn: z.string().min(1, 'Informe a data inicial.'),
  endOn: z.string().min(1, 'Informe a data final.'),
})

const periodSchema = basePeriodSchema
  .refine((value) => fromDateInputValue(value.endOn) >= fromDateInputValue(value.startOn), {
    message: 'A data final deve ser igual ou posterior a inicial.',
    path: ['endOn'],
  })

const transactionsSchema = basePeriodSchema
  .extend({
    type: z.enum(['income', 'expense']),
  })
  .refine((value) => fromDateInputValue(value.endOn) >= fromDateInputValue(value.startOn), {
    message: 'A data final deve ser igual ou posterior a inicial.',
    path: ['endOn'],
  })

const accountSchema = basePeriodSchema
  .extend({
    accountId: z.string().trim().min(1, 'Selecione uma conta.'),
  })
  .refine((value) => fromDateInputValue(value.endOn) >= fromDateInputValue(value.startOn), {
    message: 'A data final deve ser igual ou posterior a inicial.',
    path: ['endOn'],
  })

function buildPeriodFilters(form: PeriodForm) {
  const startDate = fromDateInputValue(form.startOn)
  const endDate = fromDateInputValue(form.endOn)

  startDate.setHours(0, 0, 0, 0)
  endDate.setHours(23, 59, 59, 0)

  return {
    startAt: formatLocalRFC3339(startDate),
    endAt: formatLocalRFC3339(endDate),
  }
}

function firstErrors(fieldErrors: Record<string, string[] | undefined>) {
  return Object.fromEntries(
    Object.entries(fieldErrors).flatMap(([field, messages]) => {
      const message = messages?.[0]
      return message ? [[field, message]] : []
    }),
  )
}

function ReportItemButton({
  report,
  isPending,
  onClick,
}: {
  report: ReportItem
  isPending: boolean
  onClick: () => void
}) {
  return (
    <button
      type="button"
      className="grid min-h-[48px] w-full cursor-pointer grid-cols-[28px_minmax(0,1fr)] items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors hover:bg-[#f8fafc] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#2563eb] md:px-4"
      onClick={onClick}
    >
      <span className="grid h-7 w-7 place-items-center rounded-full bg-[#edf5ff] text-[#2563eb]">
        {report.icon}
      </span>
      <span className="min-w-0">
        <span className="block truncate text-[13px] font-semibold text-[#1f2937]">
          {report.title}
        </span>
        <span className="block truncate text-[11px] font-medium text-[#9aa5b1]">
          {isPending ? 'Gerando PDF...' : report.description}
        </span>
      </span>
    </button>
  )
}

function TypeSelector({
  value,
  onChange,
}: {
  value: ReportTransactionType
  onChange: (value: ReportTransactionType) => void
}) {
  return (
    <TransactionFieldRow
      label="Tipo"
      icon={<ListFilter className="h-5 w-5" aria-hidden="true" />}
    >
      <div className="grid w-full grid-cols-2 gap-2">
        {(['income', 'expense'] as const).map((type) => (
          <button
            key={type}
            type="button"
            className={`h-10 cursor-pointer rounded-lg text-[13px] font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
              value === type ? 'bg-[#2563eb] text-white' : 'bg-[#f4f1f7] text-[#5f536d]'
            }`}
            onClick={() => onChange(type)}
          >
            {type === 'income' ? 'Receitas' : 'Despesas'}
          </button>
        ))}
      </div>
    </TransactionFieldRow>
  )
}

export function ReportDownloadPanel() {
  const [pendingReport, setPendingReport] = useState<ReportID | null>(null)
  const [openReport, setOpenReport] = useState<Exclude<ReportID, 'accounts'> | null>(null)
  const [message, setMessage] = useState('')
  const [periodForm, setPeriodForm] = useState<PeriodForm>({ startOn: today, endOn: today })
  const [transactionsForm, setTransactionsForm] = useState<TransactionsForm>({
    startOn: today,
    endOn: today,
    type: 'income',
  })
  const [accountForm, setAccountForm] = useState<AccountForm>({
    startOn: today,
    endOn: today,
    accountId: '',
  })
  const [selectedMonth, setSelectedMonth] = useState<SelectedMonth>(() => getCurrentSelectedMonth())
  const [errors, setErrors] = useState<Record<string, string>>({})

  async function runDownload(reportID: ReportID, download: () => Promise<void>) {
    setPendingReport(reportID)
    setMessage('')

    try {
      await download()
      setOpenReport(null)
      setErrors({})
    } catch {
      setMessage('Nao foi possivel gerar este relatorio agora.')
    } finally {
      setPendingReport(null)
    }
  }

  async function downloadAccountsReport() {
    await runDownload('accounts', async () => {
      const response = await downloadAccountsReportPDF()
      savePdfResponse(response, 'relatorio-contas.pdf')
    })
  }

  async function downloadTransactionsReport() {
    const parsed = transactionsSchema.safeParse(transactionsForm)
    if (!parsed.success) {
      setErrors(firstErrors(parsed.error.flatten().fieldErrors))
      return
    }

    await runDownload('transactions', async () => {
      const response = await downloadTransactionsReportPDF({
        ...buildPeriodFilters(parsed.data),
        type: parsed.data.type,
      })
      savePdfResponse(response, 'relatorio-transacoes.pdf')
    })
  }

  async function downloadPeriodReport() {
    const parsed = periodSchema.safeParse(periodForm)
    if (!parsed.success) {
      setErrors(firstErrors(parsed.error.flatten().fieldErrors))
      return
    }

    await runDownload('period', async () => {
      const response = await downloadPeriodReportPDF(buildPeriodFilters(parsed.data))
      savePdfResponse(response, 'relatorio-periodo.pdf')
    })
  }

  async function downloadMonthlyReport() {
    await runDownload('monthly', async () => {
      const response = await downloadMonthlyReportPDF(getMonthPeriod(selectedMonth))
      savePdfResponse(response, 'relatorio-mensal.pdf')
    })
  }

  async function downloadAccountReport() {
    const parsed = accountSchema.safeParse(accountForm)
    if (!parsed.success) {
      setErrors(firstErrors(parsed.error.flatten().fieldErrors))
      return
    }

    await runDownload('account', async () => {
      const response = await downloadAccountReportPDF({
        ...buildPeriodFilters(parsed.data),
        accountId: parsed.data.accountId,
      })
      savePdfResponse(response, 'relatorio-conta.pdf')
    })
  }

  function openFilteredReport(reportID: Exclude<ReportID, 'accounts'>) {
    setOpenReport(reportID)
    setErrors({})
    setMessage('')
  }

  function handleReportClick(reportID: ReportID) {
    if (reportID === 'accounts') {
      void downloadAccountsReport()
      return
    }

    openFilteredReport(reportID)
  }

  return (
    <div className="grid border-t border-[#edf1f6] bg-white px-4 py-2 md:px-8 lg:px-10">
      <div className="grid gap-1 border-l border-[#e4ebf3] pl-4 md:pl-6">
        {reportItems.map((report) => (
          <ReportItemButton
            key={report.id}
            report={report}
            isPending={pendingReport === report.id}
            onClick={() => handleReportClick(report.id)}
          />
        ))}
      </div>

      {message ? (
        <p className="mt-2 rounded-lg bg-[#f8fafc] px-3 py-2 text-[12px] font-semibold text-[#667085]">
          {message}
        </p>
      ) : null}

      <SelectionSheet
        title="Receitas/despesas por periodo"
        isOpen={openReport === 'transactions'}
        onClose={() => setOpenReport(null)}
      >
        <div className="grid gap-0">
          <DateInput
            value={transactionsForm.startOn}
            label="Inicial"
            accentColor="#2563eb"
            error={errors.startOn}
            onChange={(startOn) => setTransactionsForm((current) => ({ ...current, startOn }))}
          />
          <DateInput
            value={transactionsForm.endOn}
            label="Final"
            accentColor="#2563eb"
            error={errors.endOn}
            onChange={(endOn) => setTransactionsForm((current) => ({ ...current, endOn }))}
          />
          <TypeSelector
            value={transactionsForm.type}
            onChange={(type) => setTransactionsForm((current) => ({ ...current, type }))}
          />
        </div>
        <SheetSubmitButton
          isPending={pendingReport === 'transactions'}
          onClick={downloadTransactionsReport}
        />
      </SelectionSheet>

      <SelectionSheet
        title="Geral por periodo"
        isOpen={openReport === 'period'}
        onClose={() => setOpenReport(null)}
      >
        <div className="grid gap-0">
          <DateInput
            value={periodForm.startOn}
            label="Inicial"
            accentColor="#2563eb"
            error={errors.startOn}
            onChange={(startOn) => setPeriodForm((current) => ({ ...current, startOn }))}
          />
          <DateInput
            value={periodForm.endOn}
            label="Final"
            accentColor="#2563eb"
            error={errors.endOn}
            onChange={(endOn) => setPeriodForm((current) => ({ ...current, endOn }))}
          />
        </div>
        <SheetSubmitButton isPending={pendingReport === 'period'} onClick={downloadPeriodReport} />
      </SelectionSheet>

      <SelectionSheet
        title="Mensal consolidado"
        isOpen={openReport === 'monthly'}
        onClose={() => setOpenReport(null)}
      >
        <div className="grid justify-items-center gap-5 px-2 py-4">
          <MonthSelector selectedMonth={selectedMonth} onSelectMonth={setSelectedMonth} />
        </div>
        <SheetSubmitButton isPending={pendingReport === 'monthly'} onClick={downloadMonthlyReport} />
      </SelectionSheet>

      <SelectionSheet
        title="Por conta bancaria"
        isOpen={openReport === 'account'}
        onClose={() => setOpenReport(null)}
      >
        <div className="grid gap-0">
          <AccountSelector
            label="Conta"
            value={accountForm.accountId}
            error={errors.accountId}
            onChange={(accountId) => setAccountForm((current) => ({ ...current, accountId }))}
          />
          <DateInput
            value={accountForm.startOn}
            label="Inicial"
            accentColor="#2563eb"
            error={errors.startOn}
            onChange={(startOn) => setAccountForm((current) => ({ ...current, startOn }))}
          />
          <DateInput
            value={accountForm.endOn}
            label="Final"
            accentColor="#2563eb"
            error={errors.endOn}
            onChange={(endOn) => setAccountForm((current) => ({ ...current, endOn }))}
          />
        </div>
        <SheetSubmitButton isPending={pendingReport === 'account'} onClick={downloadAccountReport} />
      </SelectionSheet>
    </div>
  )
}

function SheetSubmitButton({
  isPending,
  onClick,
}: {
  isPending: boolean
  onClick: () => void
}) {
  return (
    <div className="sticky bottom-0 bg-white/96 px-4 pb-[calc(16px+env(safe-area-inset-bottom))] pt-5 backdrop-blur md:static md:bg-transparent md:px-5 md:pb-5 md:pt-6 md:backdrop-blur-none">
      <button
        type="button"
        disabled={isPending}
        className="mx-auto block h-12 w-full max-w-[420px] cursor-pointer rounded-lg bg-[#281d35] px-4 text-[15px] font-semibold text-white shadow-[0_6px_14px_rgba(40,29,53,0.10)] transition-colors hover:bg-[#3a2a4a] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] disabled:cursor-not-allowed disabled:opacity-65"
        onClick={onClick}
      >
        {isPending ? 'Gerando...' : 'Gerar PDF'}
      </button>
    </div>
  )
}
