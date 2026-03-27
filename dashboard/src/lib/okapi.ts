// app/lib/okapi.ts
import { env } from '@/config/env'
import type {
  ServicesResponse,
  HealthAllResponse,
  SelfHealth,
  ServiceRecord,
} from '@/types/health'
import type { IncidentsResponse } from '@/types/incidents'
import type { MaintenanceResponse } from '@/types/maintenance'

function normalizeServiceRecord(record: any): ServiceRecord {
  return {
    ...record,
    components: record.components || [],
    incidents: record.incidents || [],
    scheduled_maintenance: record.scheduled_maintenance || [],
  }
}

async function get<T>(path: string, params?: Record<string, any>): Promise<T> {
  // Ensure we don't have double slashes
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  let url = `${env.okapiBaseUrl}${cleanPath}`
  
  if (env.okapiBaseUrl.endsWith('/') && cleanPath.startsWith('/')) {
    url = `${env.okapiBaseUrl}${cleanPath.substring(1)}`
  }

  if (params) {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, val]) => {
      if (val !== undefined && val !== null) {
        searchParams.append(key, String(val))
      }
    })
    const qs = searchParams.toString()
    if (qs) {
      url += (url.includes('?') ? '&' : '?') + qs
    }
  }

  const res = await fetch(url)
  if (!res.ok) throw new Error(`Okapi ${path} → ${res.status}`)

  const contentType = res.headers.get("content-type")
  if (!contentType || !contentType.includes("application/json")) {
    throw new Error(`Okapi ${path} returned non-JSON content: ${contentType}`)
  }

  const data = await res.json()

  // Normalize results arrays in different response shapes
  if (data) {
    if (path === '/services' && !Array.isArray(data)) {
      return [] as unknown as T
    }
  }

  return data as T
}

export const okapi = {
  selfHealth: ()  => get<SelfHealth>('/_health'),
  services:   ()  => get<ServicesResponse>('/services').then(res => Array.isArray(res) ? res : []),
  healthAll:  async (params?: any): Promise<HealthAllResponse> => {
    const services = await okapi.services()
    if (!services || services.length === 0) {
      return { results: [], count: 0 }
    }
    const res = await get<any>('/health', { services: services.join(',') })
    
    const resultsMap = res?.results || {}
    const fullResults = Object.values(resultsMap).map(normalizeServiceRecord) as ServiceRecord[]

    // Calculate status counts on all results before filtering
    const statusCounts: Record<string, number> = {
      All: fullResults.length,
      Operational: fullResults.filter(r => r.status === 'operational').length,
      Degraded: fullResults.filter(r => r.status === 'degraded' || r.status === 'partial_outage').length,
      Outage: fullResults.filter(r => r.status === 'major_outage').length,
    }

    let resultsArray = [...fullResults]
    
    if (params?.q) {
      const lowerQ = params.q.toLowerCase()
      resultsArray = resultsArray.filter(r => r.service.toLowerCase().includes(lowerQ) || r.data_source.toLowerCase().includes(lowerQ))
    }
    
    if (params?.status && params.status !== 'All') {
        const targetStatus = params.status.toLowerCase()
        if (targetStatus === 'degraded') {
            resultsArray = resultsArray.filter(r => r.status === 'degraded' || r.status === 'partial_outage')
        } else if (targetStatus === 'outage') {
            resultsArray = resultsArray.filter(r => r.status === 'major_outage')
        } else {
            resultsArray = resultsArray.filter(r => r.status.toLowerCase().includes(targetStatus))
        }
    }
    
    const count = resultsArray.length
    
    if (params?.offset !== undefined || params?.limit !== undefined) {
      const offset = Number(params?.offset) || 0
      const limit = Number(params?.limit) || 20
      resultsArray = resultsArray.slice(offset, offset + limit)
    }
    
    return {
      results: resultsArray,
      count: count,
      limit: params?.limit,
      offset: params?.offset,
      statusCounts: statusCounts
    }
  },
  stats:      (service?: string) => get<any>(service ? `/stats/${service}` : '/stats'),
  incidents:  async (params?: any): Promise<IncidentsResponse> => {
    const health = await okapi.healthAll({ limit: 1000, ...params })
    const affected = health.results.filter(r => r.incidents && r.incidents.length > 0)
    return {
      incidents: affected,
      count: affected.length
    }
  },
  maintenance: async (params?: any): Promise<MaintenanceResponse> => {
    const health = await okapi.healthAll({ limit: 1000, ...params })
    const affected = health.results.filter(r => r.scheduled_maintenance && r.scheduled_maintenance.length > 0)
    return {
      scheduled_maintenance: affected,
      count: affected.length
    }
  },
} as const
