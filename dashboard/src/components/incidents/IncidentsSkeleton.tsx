import { Skeleton } from '@/components/ui/skeleton'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBarSkeleton } from '@/components/shared/FilterBarSkeleton'

const STATUS_FILTERS = ['All', 'Investigating', 'Identified', 'Monitoring', 'Resolved'] as const

export function IncidentsSkeleton() {
  return (
    <div className="flex flex-col">
      <SectionLabel>Incidents</SectionLabel>

      <FilterBarSkeleton statuses={STATUS_FILTERS} />

      <section>
        <div className="flex flex-col border-t border-[var(--border)]">
          {['i1', 'i2', 'i3', 'i4', 'i5', 'i6'].map((key) => (
            <div key={key} className="flex items-center justify-between py-4 px-4 border-b border-[var(--border)]">
              <div className="flex items-center gap-4 w-full text-left">
                <Skeleton className="h-5 w-24 rounded-[2px] bg-[var(--surface-3)]" />
                <div className="flex flex-col gap-1.5 flex-1">
                  <Skeleton className="h-4 w-32 bg-[var(--surface-3)]" />
                  <Skeleton className="h-3 w-48 bg-[var(--surface-3)]" />
                </div>
                <div className="ml-auto flex items-center gap-6 pr-4 hidden sm:flex">
                  <Skeleton className="h-3 w-24 bg-[var(--surface-3)]" />
                  <Skeleton className="h-3 w-24 bg-[var(--surface-3)]" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
