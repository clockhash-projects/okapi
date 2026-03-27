import type { ServiceRecord } from '@/types/health'
import { UptimeGraph } from '@/components/shared/UptimeGraph'

interface ServiceCardProps {
  record: ServiceRecord
  onOpen?: (record: ServiceRecord) => void
}

function StatusDot({ status }: { status: string }) {
  const colorMap: Record<string, string> = {
    operational: 'bg-[var(--green)]',
    degraded: 'bg-[var(--yellow)]',
    partial_outage: 'bg-[var(--yellow)]',
    major_outage: 'bg-[var(--red)]',
    maintenance: 'bg-[var(--blue)]',
    unknown: 'bg-[var(--text-muted)]',
  }
  const color = colorMap[status] || colorMap.unknown
  const isOperational = status === 'operational'

  return (
    <span className="relative flex h-2.5 w-2.5 flex-shrink-0">
      {isOperational && (
        <span
          className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-60 ${color}`}
        />
      )}
      <span className={`relative inline-flex rounded-full h-2.5 w-2.5 ${color}`} />
    </span>
  )
}

function formatRelativeTime(ts: string | undefined): string {
  if (!ts) return '—'
  const diff = Date.now() - new Date(ts).getTime()
  const sec = Math.floor(diff / 1000)
  if (sec < 60) return `${sec}s ago`
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min}m ago`
  const hr = Math.floor(min / 60)
  return `${hr}h ago`
}

function calcUptime(history: any[] | undefined): string {
  if (!history || history.length === 0) return '—'
  const up = history.filter((h) => h.status === 'operational').length
  return ((up / history.length) * 100).toFixed(1) + '%'
}

export function ServiceCard({ record, onOpen }: ServiceCardProps) {
  const uptime = (record as any).uptime_str || calcUptime(record.recent_history)
  const lastChecked = formatRelativeTime(record.fetched_at as string | undefined)
  const latency = (record as any).latency_ms != null ? `${(record as any).latency_ms}ms` : null

  return (
    <div
      onClick={() => onOpen?.(record)}
      className="bg-[var(--surface)] border border-[var(--border)] px-4 py-3 rounded-[2px] flex flex-col gap-3
        hover:border-[var(--border-hover)] hover:bg-[var(--surface-2)] transition-all group cursor-pointer"
    >
      {/* Top row: dot + name + latency */}
      <div className="flex items-center gap-2.5">
        <StatusDot status={record.status} />
        <span className="font-mono text-sm text-[var(--text)] group-hover:text-white transition-colors truncate flex-1">
          {record.service}
        </span>
        {latency && (
          <span className="font-mono text-[10px] text-[var(--text-muted)] tabular-nums flex-shrink-0">
            {latency}
          </span>
        )}
      </div>

      {/* Uptime bar */}
      <UptimeGraph history={record.recent_history} limit={40} />

      {/* Bottom row: uptime % + last checked */}
      <div className="flex items-center justify-between">
        <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase tracking-wider">
          {uptime !== '—' ? `${uptime} uptime` : record.data_source}
        </span>
        <span className="font-mono text-[10px] text-[var(--text-muted)]">{lastChecked}</span>
      </div>
    </div>
  )
}
