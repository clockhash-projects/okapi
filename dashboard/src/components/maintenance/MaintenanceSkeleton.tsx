import { Skeleton } from '@/components/ui/skeleton'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBarSkeleton } from '@/components/shared/FilterBarSkeleton'

const SECTION_TABS = ['Active', 'Upcoming', 'Completed'] as const

export function MaintenanceSkeleton() {
  return (
    <div className="flex flex-col">
      <SectionLabel>Maintenance</SectionLabel>

      <FilterBarSkeleton statuses={SECTION_TABS} />

      <div className="flex flex-col border-t border-[var(--border)]">
        {['mn1', 'mn2', 'mn3', 'mn4', 'mn5', 'mn6'].map((key) => (
          <div key={key} className="flex items-center justify-between py-4 px-4 border-b border-[var(--border)]">
            <div className="flex items-center gap-4 w-full text-left">
              <Skeleton className="h-5 w-24 rounded-[2px] bg-[var(--surface-3)]" />
              <Skeleton className="h-4 w-32 bg-[var(--surface-3)]" />
              <div className="ml-auto flex items-center gap-6 pr-4 hidden sm:flex">
                <Skeleton className="h-3 w-20 bg-[var(--surface-3)]" />
                <Skeleton className="h-3 w-24 bg-[var(--surface-3)]" />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
