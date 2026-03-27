import type { StatusPoint, ServiceStatus } from '@/types/health'
import { clsx } from 'clsx'

interface UptimeGraphProps {
  history?: StatusPoint[]
  limit?: number
}

export function UptimeGraph({ history = [], limit = 40 }: UptimeGraphProps) {
  const paddedHistory = [...history].reverse()
  
  const getStatusColor = (status: ServiceStatus) => {
    switch (status) {
      case 'operational': return 'bg-[var(--green)]'
      case 'degraded':
      case 'partial_outage': return 'bg-[var(--yellow)]'
      case 'major_outage': return 'bg-[var(--red)]'
      default: return 'bg-[var(--border)]'
    }
  }

  const formatTime = (isoString: string) => {
    return new Date(isoString).toLocaleString([], { 
      month: 'short', 
      day: '2-digit', 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }

  return (
    <div className="flex gap-[2px] h-6 w-full items-end group/graph">
      {paddedHistory.map((point, i) => (
        <div
          key={i}
          className={clsx(
            'flex-1 h-full rounded-[1px] transition-all hover:scale-y-125 hover:opacity-100 opacity-80',
            getStatusColor(point.status)
          )}
          title={`${point.status} @ ${formatTime(point.time)}`}
        />
      ))}
      {paddedHistory.length === 0 && Array.from({ length: limit }).map((_, i) => (
        <div key={i} className="flex-1 h-full rounded-[1px] bg-[var(--border)] opacity-20" />
      ))}
    </div>
  )
}
