
interface MetricTileProps {
  label: string
  value: number | string
  color?: string
}

export function MetricTile({ label, value, color }: MetricTileProps) {
  return (
    <div className="bg-[var(--surface)] border border-[var(--border)] p-4 rounded-[2px] flex flex-col gap-1">
      <span className="font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">
        / {label}
      </span>
      <span
        className="font-mono text-2xl font-bold"
        style={color ? { color: `var(--${color})` } : undefined}
      >
        {value}
      </span>
    </div>
  )
}
