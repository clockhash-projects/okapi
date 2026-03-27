import { Search, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { ServiceRecord, Component } from '@/types/health'
import type { Incident } from '@/types/incidents'
import type { MaintenanceWindow } from '@/types/maintenance'
import { StatusPill } from '@/components/shared/StatusPill'

interface ServiceDrawerProps {
  record: ServiceRecord | null
  onClose: () => void
}

export function ServiceDrawer({ record, onClose }: ServiceDrawerProps) {
  const [compSearch, setCompSearch] = useState('')

  useEffect(() => {
    if (record) {
      document.body.style.overflow = 'hidden'
      setCompSearch('') // Reset search when opening new record
    } else {
      document.body.style.overflow = ''
    }
    return () => {
      document.body.style.overflow = ''
    }
  }, [record])

  if (!record) return null

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/60 backdrop-blur-sm transition-opacity animate-in fade-in duration-300" 
        onClick={onClose}
      />
      
      {/* Drawer */}
      <div className="relative w-full max-w-xl bg-[var(--surface)] border-l border-[var(--border)] h-full shadow-2xl flex flex-col animate-in slide-in-from-right duration-300 ease-out">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-[var(--border)] bg-[var(--surface-2)]">
          <div className="flex flex-col gap-1">
            <h2 className="text-xl font-mono font-bold text-white tracking-tight">{record.service}</h2>
            <a 
              href={record.source_url} 
              target="_blank" 
              rel="noreferrer"
              className="text-[10px] font-mono text-[var(--text-muted)] hover:text-white transition-colors"
            >
              {record.source_url}
            </a>
          </div>
          <div className="flex items-center gap-4">
            <StatusPill status={record.status} />
            <button 
              onClick={onClose}
              className="p-2 hover:bg-[var(--surface)] rounded-full transition-colors text-[var(--text-muted)] hover:text-white"
            >
              <X size={20} />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-8 custom-scrollbar pb-12">
          {/* Status Section */}
          <section className="flex flex-col gap-4">
            <div className="flex items-center justify-between">
              <h3 className="text-xs font-mono font-bold text-[var(--text-muted)] uppercase tracking-widest">Status Details</h3>
              <span className="text-[10px] font-mono text-[var(--text-muted)]">Source: {record.data_source}</span>
            </div>
            <div className="bg-[var(--surface-2)] p-6 border border-[var(--border)] rounded-[2px]">
              <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-col">
                  <span className="text-[10px] font-mono text-[var(--text-muted)] uppercase">Status</span>
                  <span className="text-sm font-mono text-white capitalize">{record.status.replace('_', ' ')}</span>
                </div>
                <div className="flex flex-col">
                  <span className="text-[10px] font-mono text-[var(--text-muted)] uppercase">Last Checked</span>
                  <span className="text-sm font-mono text-white">{new Date(record.fetched_at).toLocaleTimeString()}</span>
                </div>
              </div>
            </div>
          </section>

          {/* Components Section */}
          {record.components && record.components.length > 0 && (
            <section className="flex flex-col gap-4">
              <div className="flex items-center justify-between">
                <h3 className="text-xs font-mono font-bold text-[var(--text-muted)] uppercase tracking-widest">Sub-Services & Components</h3>
                <span className="text-[10px] font-mono text-[var(--text-muted)]">{record.components.length} Total</span>
              </div>

              {record.components.length > 8 && (
                <div className="relative group">
                  <Search size={12} className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)] group-focus-within:text-[var(--blue)] transition-colors" />
                  <input 
                    type="text"
                    placeholder="Search components..."
                    value={compSearch}
                    onChange={(e) => setCompSearch(e.target.value)}
                    className="w-full bg-[var(--surface-3)] border border-[var(--border)] rounded-[1px] py-2 pl-8 pr-4 text-xs font-mono text-white focus:outline-none focus:border-[var(--blue)] transition-all"
                  />
                </div>
              )}

              <div className="grid grid-cols-1 gap-2">
                {record.components
                  .filter((c: Component) => c.name.toLowerCase().includes(compSearch.toLowerCase()))
                  .map((comp: Component) => (
                    <div 
                      key={comp.name} 
                      className="flex items-center justify-between p-3 bg-[var(--surface-2)] border border-[var(--border)] rounded-[1px] hover:border-[var(--border-hover)] transition-all"
                    >
                      <span className="text-sm font-mono text-[var(--text)]">{comp.name}</span>
                      <StatusPill status={comp.status} />
                    </div>
                  ))}
              </div>
            </section>
          )}

          {/* Incidents Section */}
          {record.incidents && record.incidents.length > 0 && (
            <section className="flex flex-col gap-4">
              <h3 className="text-xs font-mono font-bold text-[var(--text-muted)] uppercase tracking-widest text-[#ff4d4d]">Active & Recent Incidents</h3>
              <div className="flex flex-col gap-4">
                {record.incidents.map((incident: Incident) => (
                  <div key={incident.id} className="p-4 bg-[var(--surface-2)] border-l-2 border-[#ff4d4d] flex flex-col gap-2">
                    <div className="flex justify-between items-start gap-4">
                      <h4 className="text-sm font-mono font-bold text-white leading-tight">{incident.title}</h4>
                      <span className="text-[10px] font-mono text-[var(--text-muted)] whitespace-nowrap">
                        {new Date(incident.created_at).toLocaleString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                      </span>
                    </div>
                    <p className="text-xs font-mono text-[var(--text-muted)] leading-relaxed">{incident.body}</p>
                    <div className="mt-1">
                      <span className="text-[9px] font-mono bg-red-500/10 text-red-400 px-1.5 py-0.5 rounded-[1px] uppercase border border-red-500/20">
                        {incident.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* Maintenance Section */}
          {record.scheduled_maintenance && record.scheduled_maintenance.length > 0 && (
            <section className="flex flex-col gap-4">
              <h3 className="text-xs font-mono font-bold text-[var(--text-muted)] uppercase tracking-widest text-blue-400">Scheduled Maintenance</h3>
              <div className="flex flex-col gap-4">
                {record.scheduled_maintenance.map((m: MaintenanceWindow) => (
                  <div key={m.id} className="p-4 bg-[var(--surface-2)] border-l-2 border-blue-400 flex flex-col gap-2">
                    <h4 className="text-sm font-mono font-bold text-white">{m.title}</h4>
                    <p className="text-xs font-mono text-[var(--text-muted)]">{m.summary}</p>
                    <div className="flex justify-between mt-2 pt-2 border-t border-[var(--border)]">
                      <div className="flex flex-col">
                        <span className="text-[8px] font-mono text-[var(--text-muted)] uppercase">Starts</span>
                        <span className="text-[10px] font-mono text-white">{new Date(m.starts_at).toLocaleString()}</span>
                      </div>
                      <div className="flex flex-col text-right">
                        <span className="text-[8px] font-mono text-[var(--text-muted)] uppercase">Ends</span>
                        <span className="text-[10px] font-mono text-white">{new Date(m.ends_at).toLocaleString()}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          )}
        </div>
      </div>
    </div>
  )
}
