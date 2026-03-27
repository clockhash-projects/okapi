
interface CodeBlockProps {
  content: unknown
}

export function CodeBlock({ content }: CodeBlockProps) {
  return (
    <pre className="bg-[var(--surface-2)] border border-[var(--border)] p-4 rounded-[2px] font-mono text-[11px] overflow-x-auto text-[var(--text-2)] leading-relaxed">
      <code>{JSON.stringify(content, null, 2)}</code>
    </pre>
  )
}
