import { env } from '@/config/env'
import { useApiHealth } from '@/hooks/useApiHealth'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { NavLinks } from './NavLinks'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function Sidebar() {
  const { reachable } = useApiHealth()

  return (
    <aside className="hidden lg:flex w-60 h-screen border-r border-[var(--border)] flex-col bg-[var(--bg)] fixed left-0 top-0 z-20">
      <div className="p-6 pb-8">
        <div className="flex items-center gap-2 mb-1">
          <div className="pulse-dot" />
          <span className="font-mono font-bold text-sm text-[var(--text)] uppercase tracking-tight">okapi</span>
        </div>
        <div className="font-mono text-[9px] text-[var(--text-muted)] tracking-[0.2em] uppercase">
          / Universal Health API
        </div>
      </div>

      <nav className="flex-1 flex flex-col">
        <NavLinks />
      </nav>

      <div className="mt-auto p-4 border-t border-[var(--border)]">
        <div className="flex items-center gap-2 mb-2">
          <div className={cn('w-1.5 h-1.5 rounded-full', reachable ? 'bg-[var(--green)]' : 'bg-[var(--red)]')} />
          <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase tracking-wider">
            {reachable ? 'connected' : 'offline'}
          </span>
        </div>
        <div className="font-mono text-[9px] text-[var(--text-muted)] truncate" title={env.okapiBaseUrl}>
          {env.okapiBaseUrl.replace(/^https?:\/\//, '')}
        </div>
      </div>
    </aside>
  )
}
