import type { ServiceStatus } from '@/types/health'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

interface StatusPillProps {
  status: ServiceStatus
  className?: string
}

const statusMap: Record<ServiceStatus, { label: string; className: string }> = {
  operational: {
    label: 'operational',
    className: 'text-[var(--green)] bg-[var(--green-dim)]',
  },
  degraded: {
    label: 'degraded',
    className: 'text-[var(--yellow)] bg-[var(--yellow-dim)]',
  },
  partial_outage: {
    label: 'partial outage',
    className: 'text-[var(--yellow)] bg-[var(--yellow-dim)]',
  },
  major_outage: {
    label: 'major outage',
    className: 'text-[var(--red)] bg-[var(--red-dim)]',
  },
  maintenance: {
    label: 'maintenance',
    className: 'text-[var(--blue)] bg-[var(--blue-dim)]',
  },
  unknown: {
    label: 'unknown',
    className: 'text-[var(--text-muted)] bg-transparent border border-[var(--border)]',
  },
}

export function StatusPill({ status, className }: StatusPillProps) {
  const config = statusMap[status] || statusMap.unknown

  return (
    <span
      className={cn(
        'font-mono text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-[2px] inline-flex items-center whitespace-nowrap',
        config.className,
        className
      )}
    >
      {config.label}
    </span>
  )
}
