import type { ServiceRecord } from '@/types/health'
import { StatusPill } from '@/components/shared/StatusPill'

interface ServiceTableProps {
  records: ServiceRecord[]
  onOpen: (record: ServiceRecord) => void
  onTogglePin?: (id: string) => void
  isPinned?: (id: string) => boolean
}

function getRelativeTime(isoString: string) {
  if (!isoString) return 'Never'
  const date = new Date(isoString)
  if (isNaN(date.getTime())) return 'Unknown'
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (diffInSeconds < 60) return `${diffInSeconds}s ago`
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`
  return date.toLocaleDateString()
}

function PinButton({
  pinned,
  onClick,
}: {
  pinned: boolean
  onClick: (e: React.MouseEvent) => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      title={pinned ? 'Unpin service' : 'Pin service'}
      className={`p-1.5 rounded transition-colors ${
        pinned
          ? 'text-[var(--text)] opacity-100'
          : 'text-[var(--text-muted)] opacity-20 group-hover:opacity-60 hover:!opacity-100'
      }`}
    >
      <svg
        width="13"
        height="13"
        viewBox="0 0 24 24"
        fill={pinned ? 'currentColor' : 'none'}
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <line x1="12" y1="17" x2="12" y2="22" />
        <path d="M5 17h14v-1.76a2 2 0 0 0-1.11-1.79l-1.78-.9A2 2 0 0 1 15 10.76V6h1a2 2 0 0 0 0-4H8a2 2 0 0 0 0 4h1v4.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24Z" />
      </svg>
    </button>
  )
}

export function ServiceTable({ records, onOpen, onTogglePin, isPinned }: ServiceTableProps) {
  return (
    <div className="overflow-x-auto border-t border-[var(--border)]">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-[var(--border)]">
            {onTogglePin && (
              <th className="py-3 px-3 w-8">
                <div className="flex justify-center text-[var(--text-muted)] opacity-50">
                  <svg
                    width="10"
                    height="10"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <line x1="12" y1="17" x2="12" y2="22" />
                    <path d="M5 17h14v-1.76a2 2 0 0 0-1.11-1.79l-1.78-.9A2 2 0 0 1 15 10.76V6h1a2 2 0 0 0 0-4H8a2 2 0 0 0 0 4h1v4.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24Z" />
                  </svg>
                </div>
              </th>
            )}
            <th className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">
              Service
            </th>
            <th className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">Status</th>
            <th className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">Summary</th>
            <th className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">Source</th>
            <th className="py-3 px-4 font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">Checked</th>
          </tr>
        </thead>
        <tbody>
          {records.map((record) => {
            const pinned = isPinned?.(record.service) ?? false
            return (
              <tr
                key={record.service}
                onClick={() => onOpen(record)}
                className={`border-b border-[var(--border)] hover:bg-[var(--surface-2)] transition-colors group cursor-pointer ${
                  pinned ? 'bg-[var(--surface)]' : ''
                }`}
              >
                {onTogglePin && (
                  <td className="py-3 px-3">
                    <PinButton
                      pinned={pinned}
                      onClick={(e) => {
                        e.stopPropagation()
                        onTogglePin(record.service)
                      }}
                    />
                  </td>
                )}
                <td className="py-3 px-4 font-mono text-sm text-[var(--text)] group-hover:text-white">
                  {record.service}
                </td>
                <td className="py-3 px-4">
                  <StatusPill status={record.status} />
                </td>
                <td className="py-3 px-4 font-mono text-sm text-[var(--text-muted)] group-hover:text-[var(--text-2)] max-w-xs truncate">
                  {record.summary}
                </td>
                <td className="py-3 px-4 font-mono text-sm text-[var(--text-muted)] group-hover:text-[var(--text-2)] uppercase">
                  {record.data_source}
                </td>
                <td className="py-3 px-4 font-mono text-sm text-[var(--text-muted)] group-hover:text-[var(--text-2)]">
                  {getRelativeTime(record.fetched_at)}
                </td>
              </tr>
            )
          })}
          {records.length === 0 && (
            <tr>
              <td colSpan={onTogglePin ? 6 : 5} className="py-12 text-center font-mono text-sm text-[var(--text-muted)]">
                No services found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}
