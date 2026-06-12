import type { AxiosResponse } from 'axios'
import { api } from '../../../lib/api/axios.ts'
import {
  financialReportFiltersSchema,
  financialReportSchema,
} from '../schemas/reports.ts'
import type { FinancialReport, FinancialReportFilters } from '../types/reports.ts'

type PdfResponse = AxiosResponse<Blob>

const pdfHeaders = {
  Accept: 'application/pdf',
}

function cleanFilters(filters: FinancialReportFilters) {
  const parsed = financialReportFiltersSchema.parse(filters)

  return {
    ...parsed,
    categoryId: parsed.categoryId || undefined,
    accountId: parsed.accountId || undefined,
  }
}

export async function getFinancialReport(
  filters: FinancialReportFilters,
): Promise<FinancialReport> {
  const response = await api.get<unknown>('/reports/financial', {
    params: cleanFilters(filters),
  })

  return financialReportSchema.parse(response.data)
}

export async function downloadFinancialReportPDF(
  filters: FinancialReportFilters,
): Promise<PdfResponse> {
  return api.get<Blob>('/reports/financial/pdf', {
    headers: pdfHeaders,
    params: cleanFilters(filters),
    responseType: 'blob',
  })
}

function getFilenameFromContentDisposition(contentDisposition: string | undefined) {
  if (!contentDisposition) {
    return undefined
  }

  const encodedFilename = contentDisposition.match(/filename\*=UTF-8''([^;]+)/i)?.[1]

  if (encodedFilename) {
    return decodeURIComponent(encodedFilename)
  }

  return contentDisposition.match(/filename="?([^";]+)"?/i)?.[1]
}

export function savePdfResponse(response: PdfResponse, fallbackFilename: string) {
  const filename =
    getFilenameFromContentDisposition(response.headers['content-disposition']) ?? fallbackFilename
  const url = window.URL.createObjectURL(response.data)
  const link = document.createElement('a')

  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(url)
}
