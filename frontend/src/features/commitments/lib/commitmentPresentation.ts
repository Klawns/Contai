import {
  ArrowDownToLine,
  ArrowUpFromLine,
  type LucideIcon,
} from 'lucide-react'
import type {
  CommitmentType,
  EffectiveCommitmentStatus,
  RecurrenceFrequency,
} from '../types/commitments.ts'

export const typeCopy: Record<
  CommitmentType,
  {
    title: string
    plural: string
    newTitle: string
    editTitle: string
    settleTitle: string
    settleButton: string
    settledLabel: string
    accent: string
    header: string
    tone: 'income' | 'expense'
    icon: LucideIcon
  }
> = {
  payable: {
    title: 'Conta a pagar',
    plural: 'A pagar',
    newTitle: 'Nova conta a pagar',
    editTitle: 'Editar conta a pagar',
    settleTitle: 'Quitar conta a pagar',
    settleButton: 'Pagar',
    settledLabel: 'Paga',
    accent: '#d93658',
    header: 'bg-[#d93658]',
    tone: 'expense',
    icon: ArrowUpFromLine,
  },
  receivable: {
    title: 'Conta a receber',
    plural: 'A receber',
    newTitle: 'Nova conta a receber',
    editTitle: 'Editar conta a receber',
    settleTitle: 'Quitar conta a receber',
    settleButton: 'Receber',
    settledLabel: 'Recebida',
    accent: '#159c57',
    header: 'bg-[#159c57]',
    tone: 'income',
    icon: ArrowDownToLine,
  },
}

export const statusCopy: Record<EffectiveCommitmentStatus, string> = {
  pending: 'Pendente',
  overdue: 'Vencida',
  paid: 'Paga',
  received: 'Recebida',
  canceled: 'Cancelada',
}

export const recurrenceFrequencyCopy: Record<RecurrenceFrequency, string> = {
  daily: 'Diaria',
  weekly: 'Semanal',
  monthly: 'Mensal',
  yearly: 'Anual',
}

type RecurrenceFrequencyOption = {
  value: RecurrenceFrequency
  label: string
}

export const recurrenceFrequencyOptions: RecurrenceFrequencyOption[] = [
  { value: 'daily', label: recurrenceFrequencyCopy.daily },
  { value: 'weekly', label: recurrenceFrequencyCopy.weekly },
  { value: 'monthly', label: recurrenceFrequencyCopy.monthly },
  { value: 'yearly', label: recurrenceFrequencyCopy.yearly },
]
