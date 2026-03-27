import { Skeleton } from '@/components/ui/skeleton'

export function FilterBarSkeleton({ statuses }: { statuses: readonly string[] | string[] }) {
  return (
    <div className="flex flex-col gap-6 mb-8">
      <div className="flex flex-col gap-2">
        <div className="h-3 w-32 bg-[var(--surface-3)] opacity-20" />
        <Skeleton className="h-10 w-full max-w-md bg-[var(--surface-2)] border border-[var(--border)] rounded-[2px]" />
      </div>

      <div className="flex items-center gap-6 overflow-x-auto pb-2 scrollbar-hide">
        {statuses.map((s) => (
          <div key={s} className="flex items-center gap-2">
            <Skeleton className="h-4 w-20 bg-[var(--surface-3)]" />
            <Skeleton className="h-4 w-6 rounded-full bg-[var(--surface-3)]" />
          </div>
        ))}
      </div>
    </div>
  )
}
