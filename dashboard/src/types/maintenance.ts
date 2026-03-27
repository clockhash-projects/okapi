import type { ServiceRecord } from './health'

export type MaintenanceStatus =
  | 'scheduled'
  | 'in_progress'
  | 'completed'

export interface MaintenanceWindow {
  id: string
  title: string
  status: MaintenanceStatus
  summary: string
  starts_at: string
  ends_at: string
  updated_at: string
  [key: string]: unknown
}

export interface MaintenanceResponse {
  scheduled_maintenance: ServiceRecord[]
  count: number
  limit?: number
  offset?: number
}
