import { ServiceRecord } from './health'

export type IncidentStatus =
  | 'investigating'
  | 'identified'
  | 'monitoring'
  | 'resolved'

export interface Incident {
  id: string
  title: string
  status: IncidentStatus
  body: string
  created_at: string
  updated_at: string
}

export interface IncidentsResponse {
  incidents: ServiceRecord[]
  count: number
}
