import type { ServiceRecord } from '@/types/health'
import { StatusPill } from '@/components/shared/StatusPill'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'

interface IncidentAccordionProps {
  records: ServiceRecord[]
}

export function IncidentAccordion({ records }: IncidentAccordionProps) {
  return (
    <Accordion multiple className="w-full">
      {records.map((record) => (
        <AccordionItem
          key={record.service}
          value={record.service}
          className="border-b border-[var(--border)]"
        >
          <AccordionTrigger className="hover:no-underline hover:bg-[var(--surface-2)] px-4 py-4 transition-colors">
            <div className="flex items-center gap-4 w-full text-left">
              <StatusPill status={record.status} />
              <div className="flex flex-col gap-0.5">
                <span className="font-mono text-sm text-[var(--text)]">{record.service}</span>
                <span className="font-mono text-[10px] text-[var(--text-muted)] truncate max-w-sm">
                  {record.summary}
                </span>
              </div>
              <div className="ml-auto flex items-center gap-6 pr-4">
                <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase">
                  {record.fetched_at ? new Date(record.fetched_at).toLocaleDateString() : 'UNKNOWN'}
                </span>
                <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase">
                  {(record.components || []).length} COMPONENTS
                </span>
              </div>
            </div>
          </AccordionTrigger>
          <AccordionContent className="bg-[var(--surface)] p-6 flex flex-col gap-8">
            {/* Components Sub-table */}
            <div className="flex flex-col gap-3">
              <span className="font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">/ Components</span>
              <div className="border border-[var(--border)] rounded-[2px] overflow-hidden">
                <table className="w-full text-left">
                  <thead>
                    <tr className="bg-[var(--surface-2)] border-b border-[var(--border)]">
                      <th className="py-2 px-3 font-mono text-[9px] text-[var(--text-muted)] uppercase">Name</th>
                      <th className="py-2 px-3 font-mono text-[9px] text-[var(--text-muted)] uppercase">Status</th>
                      <th className="py-2 px-3 font-mono text-[9px] text-[var(--text-muted)] uppercase">Updated</th>
                    </tr>
                  </thead>
                  <tbody>
                    {(record.components || []).map((comp) => (
                      <tr key={comp.name} className="border-b border-[var(--border)] last:border-0">
                        <td className="py-2 px-3 font-mono text-xs">{comp.name}</td>
                        <td className="py-2 px-3">
                          <StatusPill status={comp.status} className="scale-90 origin-left" />
                        </td>
                        <td className="py-2 px-3 font-mono text-[10px] text-[var(--text-muted)]">
                          {comp.updated_at ? new Date(comp.updated_at).toLocaleString() : 'UNKNOWN'}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Incident Entries */}
            <div className="flex flex-col gap-3">
              <span className="font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">/ Incidents</span>
              <div className="flex flex-col gap-4">
                {(record.incidents || []).map((inc) => (
                  <div key={inc.id} className="bg-[var(--surface-2)] border border-[var(--border)] p-4 rounded-[2px] flex flex-col gap-3">
                    <div className="flex items-center justify-between">
                      <span className="font-mono text-sm font-bold">{inc.title}</span>
                      <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase">{inc.status}</span>
                    </div>
                    <p className="font-mono text-xs text-[var(--text-2)] leading-relaxed">
                      {inc.body}
                    </p>
                    <div className="flex items-center gap-4 mt-1">
                      <span className="font-mono text-[9px] text-[var(--text-muted)] uppercase">ID: {inc.id}</span>
                      <span className="font-mono text-[9px] text-[var(--text-muted)] uppercase">Created: {inc.created_at ? new Date(inc.created_at).toLocaleString() : 'UNKNOWN'}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  )
}
