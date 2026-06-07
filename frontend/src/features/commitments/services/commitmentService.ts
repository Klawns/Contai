import { api } from '../../../lib/api/axios.ts'
import {
  commitmentFiltersSchema,
  commitmentPayloadSchema,
  commitmentSchema,
  commitmentsSchema,
  settlementPayloadSchema,
} from '../schemas/commitments.ts'
import type {
  Commitment,
  CommitmentFilters,
  CommitmentPayload,
  CommitmentType,
  SettlementPayload,
} from '../types/commitments.ts'

const endpointByType = {
  payable: '/payables',
  receivable: '/receivables',
} satisfies Record<CommitmentType, string>

export async function listCommitments(
  type: CommitmentType,
  filters: CommitmentFilters,
): Promise<Commitment[]> {
  const params = commitmentFiltersSchema.parse(filters)
  const response = await api.get<unknown>(endpointByType[type], { params })

  return commitmentsSchema.parse(response.data)
}

export async function createCommitment(
  type: CommitmentType,
  payload: CommitmentPayload,
): Promise<Commitment> {
  const body = commitmentPayloadSchema.parse(payload)
  const response = await api.post<unknown>(endpointByType[type], body)

  return commitmentSchema.parse(response.data)
}

export async function updateCommitment(
  type: CommitmentType,
  commitmentId: string,
  payload: CommitmentPayload,
): Promise<Commitment> {
  const body = commitmentPayloadSchema.parse(payload)
  const response = await api.patch<unknown>(`${endpointByType[type]}/${commitmentId}`, body)

  return commitmentSchema.parse(response.data)
}

export async function settleCommitment(
  type: CommitmentType,
  commitmentId: string,
  payload: SettlementPayload,
): Promise<Commitment> {
  const body = settlementPayloadSchema.parse(payload)
  const action = type === 'payable' ? 'pay' : 'receive'
  const response = await api.patch<unknown>(
    `${endpointByType[type]}/${commitmentId}/${action}`,
    body,
  )

  return commitmentSchema.parse(response.data)
}

export async function cancelCommitment(
  type: CommitmentType,
  commitmentId: string,
): Promise<Commitment> {
  const response = await api.patch<unknown>(`${endpointByType[type]}/${commitmentId}/cancel`)

  return commitmentSchema.parse(response.data)
}
