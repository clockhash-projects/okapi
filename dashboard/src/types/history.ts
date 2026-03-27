import type { ServiceRecord } from './health'

export interface HistoryResponse {
  history: ServiceRecord[]
  count: number
  limit?: number
  offset?: number
}
