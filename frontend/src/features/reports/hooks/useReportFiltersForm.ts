import { useState } from 'react'
import type { FormEvent } from 'react'
import {
  downloadFinancialReportPDF,
  savePdfResponse,
} from '../services/reportService.ts'
import type { ReportFormState } from '../types/reportForm.ts'
import {
  buildFinancialReportFilters,
  getDefaultReportFormState,
} from '../utils/reportFilters.ts'

export function useReportFiltersForm() {
  const [formState, setFormState] = useState(getDefaultReportFormState)
  const [isExporting, setIsExporting] = useState(false)
  const [message, setMessage] = useState('')

  function updateField<Key extends keyof ReportFormState>(key: Key, value: ReportFormState[Key]) {
    setFormState((current) => ({ ...current, [key]: value }))
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setIsExporting(true)
    setMessage('')

    try {
      const response = await downloadFinancialReportPDF(buildFinancialReportFilters(formState))
      savePdfResponse(response, 'relatorio-financeiro.pdf')
    } catch {
      setMessage('Nao foi possivel gerar o PDF agora.')
    } finally {
      setIsExporting(false)
    }
  }

  return {
    formState,
    updateField,
    isExporting,
    message,
    handleSubmit,
  }
}
