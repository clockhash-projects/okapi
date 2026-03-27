import type { MaintenanceWindow } from '@/types/maintenance'
import { StatusPill } from '@/components/shared/StatusPill'

interface MaintenanceListProps {
  windows: MaintenanceWindow[]
}

export function MaintenanceList({ windows }: MaintenanceListProps) {
  return (
    <div className="flex flex-col gap-4">
      {windows.map((window) => (
        <div 
          key={window.id} 
          className="bg-[var(--surface-2)] border border-[var(--border)] p-4 rounded-[2px] flex flex-col gap-3"
        >
          <div className="flex items-center justify-between">
            <span className="font-mono text-sm font-bold">{window.title}</span>
            <StatusPill status={window.status as any} />
          </div>
          
          <p className="font-mono text-xs text-[var(--text-2)] leading-relaxed">
            {window.summary}
          </p>
          
          <div className="flex flex-wrap items-center gap-x-6 gap-y-2 mt-1 pt-3 border-t border-[var(--border)]">
            <div className="flex flex-col gap-0.5">
              <span className="font-mono text-[9px] text-[var(--text-muted)] uppercase tracking-wider">Starts</span>
              <span className="font-mono text-[10px]">{new Date(window.starts_at).toLocaleString()}</span>
            </div>
            <div className="flex flex-col gap-0.5">
              <span className="font-mono text-[9px] text-[var(--text-muted)] uppercase tracking-wider">Ends</span>
              <span className="font-mono text-[10px]">{new Date(window.ends_at).toLocaleString()}</span>
            </div>
            <div className="flex flex-col gap-0.5 ml-auto">
              <span className="font-mono text-[9px] text-[var(--text-muted)] uppercase tracking-wider text-right">ID</span>
              <span className="font-mono text-[10px] text-[var(--text-muted)]">{window.id}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
