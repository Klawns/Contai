import type { AxiosResponse } from 'axios'
import { api } from '../../../lib/api/axios.ts'
import type {
  AccountReportFilters,
  ReportPeriodFilters,
  TransactionsReportFilters,
} from '../types/reports.ts'

type PdfResponse = AxiosResponse<Blob>

const pdfHeaders = {
  Accept: 'application/pdf',
}

export async function downloadAccountsReportPDF(): Promise<PdfResponse> {
  return api.get<Blob>('/reports/accounts/pdf', {
    headers: pdfHeaders,
    responseType: 'blob',
  })
}

export async function downloadTransactionsReportPDF(
  filters: TransactionsReportFilters,
): Promise<PdfResponse> {
  return api.get<Blob>('/reports/transactions/pdf', {
    headers: pdfHeaders,
    params: filters,
    responseType: 'blob',
  })
}

export async function downloadPeriodReportPDF(filters: ReportPeriodFilters): Promise<PdfResponse> {
  return api.get<Blob>('/reports/period/pdf', {
    headers: pdfHeaders,
    params: filters,
    responseType: 'blob',
  })
}

export async function downloadMonthlyReportPDF(filters: ReportPeriodFilters): Promise<PdfResponse> {
  return api.get<Blob>('/reports/monthly/pdf', {
    headers: pdfHeaders,
    params: filters,
    responseType: 'blob',
  })
}

export async function downloadAccountReportPDF(
  filters: AccountReportFilters,
): Promise<PdfResponse> {
  const { accountId, ...params } = filters

  return api.get<Blob>(`/reports/account/${accountId}/pdf`, {
    headers: pdfHeaders,
    params,
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
