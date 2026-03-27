import type { Incident } from './incidents'
import type { MaintenanceWindow } from './maintenance'

export type ServiceStatus =
  | 'operational'
  | 'degraded'
  | 'partial_outage'
  | 'major_outage'
  | 'maintenance'
  | 'unknown'

export interface StatusPoint {
  status: ServiceStatus
  time: string
}

export interface Component {
  name: string
  status: ServiceStatus
  updated_at: string
}

export interface ServiceRecord {
  service: string
  status: ServiceStatus
  summary: string
  components: Component[]
  incidents: Incident[]
  fetched_at: string
  data_source: string
  source_url: string
  cached: boolean
  scheduled_maintenance: MaintenanceWindow[]
  recent_history?: StatusPoint[]
}

export type ServicesResponse = string[]

export interface HealthAllResponse {
  results: ServiceRecord[]
  count: number
  limit?: number
  offset?: number
  statusCounts?: Record<string, number>
}

export interface SelfHealth {
  status: string
  version?: string
  uptime?: string
}
