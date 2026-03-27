import { env } from '@/config/env'
import { RefreshCw, Menu } from 'lucide-react'
import { useRouter } from '@tanstack/react-router'
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { NavLinks } from './NavLinks'

interface HeaderProps {
  title: string
  lastFetched?: string
}

export function Header({ title, lastFetched }: HeaderProps) {
  const router = useRouter()

  const handleRefresh = () => {
    router.invalidate()
  }

  return (
    <header className="h-12 border-b border-[var(--border)] flex items-center justify-between px-6 bg-[var(--bg)] sticky top-0 z-10">
      <div className="flex items-center gap-4">
        <Sheet>
          <SheetTrigger className="lg:hidden text-[var(--text-muted)] hover:text-[var(--text)] transition-colors p-1">
            <Menu size={18} />
          </SheetTrigger>
          <SheetContent side="left" className="bg-[var(--bg)] border-r border-[var(--border)] p-0 w-64">
            <SheetHeader className="p-6 pb-8 text-left">
              <SheetTitle className="flex items-center gap-2 mb-1">
                <div className="pulse-dot" />
                <span className="font-mono font-bold text-sm text-[var(--text)] uppercase tracking-tight">okapi</span>
              </SheetTitle>
              <div className="font-mono text-[9px] text-[var(--text-muted)] tracking-[0.2em] uppercase">
                / Universal Health API
              </div>
            </SheetHeader>
            <nav className="flex flex-col">
              <NavLinks />
            </nav>
          </SheetContent>
        </Sheet>
        <div className="font-mono text-sm font-medium tracking-tight uppercase">
          {title}
        </div>
      </div>

      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <span
            className={`font-mono text-[10px] uppercase px-1.5 py-0.5 rounded-[2px] ${
              env.environment === 'local'
                ? 'bg-[var(--yellow-dim)] text-[var(--yellow)]'
                : 'bg-[var(--blue-dim)] text-[var(--blue)]'
            }`}
          >
            {env.environment}
          </span>
          {lastFetched && (
            <span className="font-mono text-[10px] text-[var(--text-muted)]">
              LAST FETCHED: {new Date(lastFetched).toLocaleTimeString()}
            </span>
          )}
        </div>

        <button
          onClick={handleRefresh}
          className="text-[var(--text-muted)] hover:text-[var(--text)] transition-colors p-1"
          title="Refresh"
        >
          <RefreshCw size={14} />
        </button>
      </div>
    </header>
  )
}
