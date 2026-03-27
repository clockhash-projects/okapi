import { createFileRoute } from '@tanstack/react-router'
import { ServiceCard } from '@/components/dashboard/ServiceCard'
import { useApiHealth } from '@/hooks/useApiHealth'
import type { ServiceRecord } from '@/types/health'

export const Route = createFileRoute('/dashboard')({
  staticData: {
    title: 'Status',
  },
  component: DashboardWrapper,
})

function AllClearBanner() {
  return (
    <div className="flex items-center gap-3 px-5 py-3 bg-[var(--green-dim)] border border-[var(--green)] border-opacity-30 rounded-[2px] mb-6">
      <span className="relative flex h-2.5 w-2.5 flex-shrink-0">
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[var(--green)] opacity-60" />
        <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-[var(--green)]" />
      </span>
      <span className="font-mono text-sm text-[var(--green)]">All systems operational</span>
      <span className="font-mono text-[10px] text-[var(--green)] opacity-60 ml-auto uppercase tracking-widest">
        {new Date().toUTCString().replace('GMT', 'UTC')}
      </span>
    </div>
  )
}

function DashboardWrapper() {
  const { reachable, latency, uptime } = useApiHealth()

  // Generate fake history data since the proxy doesnt track its own history yet
  const recent_history = Array.from({ length: 40 }).map((_, i) => ({
    status: reachable ? 'operational' : ('unknown' as any),
    time: new Date(Date.now() - (40 - i) * 60000).toISOString(),
  }))

  const record: ServiceRecord = {
    service: 'Okapi Proxy (Primary)',
    status: reachable ? 'operational' : reachable === false ? 'major_outage' : 'unknown',
    summary: 'Universal proxy server',
    components: [],
    incidents: [],
    fetched_at: new Date().toISOString(),
    data_source: 'internal',
    source_url: '',
    cached: false,
    scheduled_maintenance: [],
    recent_history: recent_history,
  }

  // Inject real latency into the record if available
  if (latency !== null) {
    ;(record as any).latency_ms = latency
  }
  if (uptime !== null) {
    ;(record as any).uptime_str = uptime
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-6 flex-wrap mb-4">
        <span className="font-mono text-xs text-[var(--text-muted)]">
          <span className="text-[var(--text)]">1</span> proxy instance
        </span>
        {reachable ? (
          <span className="font-mono text-xs text-[var(--green)]">1 operational</span>
        ) : reachable === false ? (
          <span className="font-mono text-xs text-[var(--red)]">1 offline</span>
        ) : (
          <span className="font-mono text-xs text-[var(--text-muted)]">Connecting...</span>
        )}
      </div>

      {reachable && <AllClearBanner />}

      <div className="w-full">
        {reachable !== null && (
          <ServiceCard record={record} />
        )}
      </div>
    </div>
  )
}
