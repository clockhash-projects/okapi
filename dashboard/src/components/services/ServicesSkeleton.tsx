import { Skeleton } from '@/components/ui/skeleton'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBarSkeleton } from '@/components/shared/FilterBarSkeleton'

export function ServicesSkeleton() {
  return (
    <div className="flex flex-col">
      <SectionLabel>Services</SectionLabel>
      
      <FilterBarSkeleton statuses={['All', 'Operational', 'Degraded', 'Outage']} />

      {/* Table */}
      <div className="overflow-x-auto border-t border-[var(--border)]">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-[var(--border)]">
              {['Service', 'Status', 'Summary', 'Source', 'Checked', 'Cached'].map((header) => (
                <th key={header} className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {['t1', 't2', 't3', 't4', 't5', 't6', 't7', 't8', 't9', 't10'].map((key) => (
              <tr key={key} className="border-b border-[var(--border)]">
                <td className="py-3 px-4"><Skeleton className="h-4 w-32 bg-[var(--surface-3)]" /></td>
                <td className="py-3 px-4"><Skeleton className="h-5 w-24 rounded-[2px] bg-[var(--surface-3)]" /></td>
                <td className="py-3 px-4"><Skeleton className="h-4 w-64 bg-[var(--surface-3)]" /></td>
                <td className="py-3 px-4"><Skeleton className="h-4 w-20 bg-[var(--surface-3)]" /></td>
                <td className="py-3 px-4"><Skeleton className="h-4 w-16 bg-[var(--surface-3)]" /></td>
                <td className="py-3 px-4"><Skeleton className="h-4 w-8 bg-[var(--surface-3)]" /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
