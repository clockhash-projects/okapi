import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

interface FilterBarProps<T extends string> {
  search: string
  onSearchChange: (val: string) => void
  status: T
  onStatusChange: (val: T) => void
  statuses: readonly T[] | T[]
  counts?: Record<string, number>
  placeholder?: string
  label?: string
}

export function FilterBar<T extends string>({
  search,
  onSearchChange,
  status,
  onStatusChange,
  statuses,
  counts,
  placeholder = 'FILTER...',
}: FilterBarProps<T>) {
  return (
    <div className="flex flex-col gap-6 mb-8">
      <div className="flex flex-col gap-2">
        <input
          type="text"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={placeholder}
          className="bg-[var(--surface-2)] border border-[var(--border)] px-4 py-2 font-mono text-sm text-[var(--text)] focus:outline-none focus:border-[var(--border-hover)] w-full max-w-md rounded-[2px] placeholder:text-[var(--text-muted)]"
        />
      </div>

      <div className="flex items-center gap-6 overflow-x-auto pb-2 scrollbar-hide">
        {statuses.map((s) => (
          <button
            key={s}
            onClick={() => onStatusChange(s)}
            className={cn(
              'font-mono text-xs uppercase tracking-widest transition-colors whitespace-nowrap flex items-center gap-2',
              status === s ? 'text-[var(--text)]' : 'text-[var(--text-muted)] hover:text-[var(--text-2)]'
            )}
          >
            <span>{s}</span>
            {counts && typeof counts[s] !== 'undefined' && (
              <span className={cn(
                'text-[10px] px-1.5 py-0.5 rounded-full border',
                status === s
                  ? 'bg-[var(--surface-3)] border-[var(--border-hover)] text-[var(--text)]'
                  : 'bg-transparent border-[var(--border)] text-[var(--text-muted)]'
              )}>
                {counts[s]}
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  )
}
