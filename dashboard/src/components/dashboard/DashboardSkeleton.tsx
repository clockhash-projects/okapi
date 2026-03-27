export function DashboardSkeleton() {
  return (
    <div className="flex flex-col gap-6 animate-pulse">
      {/* Summary strip */}
      <div className="flex items-center gap-6">
        <div className="h-3 w-24 bg-[var(--surface-2)] rounded" />
        <div className="h-3 w-20 bg-[var(--surface-2)] rounded" />
        <div className="h-3 w-18 bg-[var(--surface-2)] rounded" />
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3">
        {Array.from({ length: 20 }).map((_, i) => (
          <div
            key={i}
            className="bg-[var(--surface)] border border-[var(--border)] px-4 py-3 rounded-[2px] flex flex-col gap-3"
          >
            {/* Top row */}
            <div className="flex items-center gap-2.5">
              <div className="h-2.5 w-2.5 rounded-full bg-[var(--border)] flex-shrink-0" />
              <div className="h-3 flex-1 bg-[var(--surface-2)] rounded" />
              <div className="h-3 w-10 bg-[var(--surface-2)] rounded flex-shrink-0" />
            </div>
            {/* Uptime bar */}
            <div className="h-4 w-full bg-[var(--surface-2)] rounded" />
            {/* Bottom row */}
            <div className="flex items-center justify-between">
              <div className="h-2 w-20 bg-[var(--border)] rounded" />
              <div className="h-2 w-12 bg-[var(--border)] rounded" />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
