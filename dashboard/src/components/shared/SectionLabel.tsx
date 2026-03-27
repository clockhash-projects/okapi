
interface SectionLabelProps {
  children: React.ReactNode
  count?: number
}

export function SectionLabel({ children, count }: SectionLabelProps) {
  return (
    <div className="flex items-center gap-2 mb-4">
      <span className="font-mono text-[10px] text-[var(--text-muted)] tracking-[0.12em] uppercase whitespace-nowrap">
        / {children} {count !== undefined && `(${count})`}
      </span>
      <div className="h-[1px] w-full bg-[var(--border)]" />
    </div>
  )
}
