import { Link } from '@tanstack/react-router'
import { FileQuestion, MoveLeft } from 'lucide-react'

export function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-8 p-4">
      {/* Glitchy/Terminal 404 */}
      <div className="relative group">
        <h1 className="font-mono text-[120px] leading-none font-black text-[var(--surface-3)] select-none group-hover:text-[var(--text)] transition-colors duration-500">
          404
        </h1>
        <div className="absolute top-0 left-0 w-full h-full flex items-center justify-center">
          <FileQuestion size={48} className="text-[var(--text-muted)] opacity-20 group-hover:opacity-100 transition-opacity" />
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 max-w-sm text-center">
        <span className="font-mono text-xs text-[var(--text-muted)] tracking-widest uppercase">
          [ ERROR_PAGE_NOT_FOUND ]
        </span>
        <p className="font-mono text-sm text-[var(--text)] leading-relaxed">
          The requested service health record or page does not exist in our federation.
        </p>
      </div>

      <Link
        to="/"
        className="flex items-center gap-2 mt-4 px-6 py-3 bg-[var(--text)] text-[var(--bg)] font-mono text-xs uppercase tracking-widest hover:opacity-90 transition-all rounded-[2px]"
      >
        <MoveLeft size={14} />
        Return to Dashboard
      </Link>
      
      {/* Visual background noise element */}
      <div className="fixed bottom-0 right-0 p-8 opacity-5 pointer-events-none select-none">
        <pre className="font-mono text-[10px] leading-tight">
          {`OKAPI STATUS: 404
PATH: UNKNOWN
ORIGIN: ${window.location.host}
TIMESTAMP: ${new_Date().toISOString()}`}
        </pre>
      </div>
    </div>
  )
}

function new_Date() {
    return new Date();
}
