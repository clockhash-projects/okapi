import { Link } from '@tanstack/react-router'

const navItems = [
  { label: 'Status', to: '/dashboard' },
  { label: 'Services', to: '/services' },
  { label: 'Incidents', to: '/incidents' },
  { label: 'Maintenance', to: '/maintenance' },
]

interface NavLinksProps {
  onLinkClick?: () => void
}

export function NavLinks({ onLinkClick }: NavLinksProps) {
  return (
    <>
      {navItems.map((item) => (
        <Link
          key={item.to}
          to={item.to}
          onClick={onLinkClick}
          className="px-6 py-2.5 font-mono text-sm text-[var(--text-muted)] hover:text-[var(--text)] transition-colors border-l-[3px] border-transparent"
          activeProps={{
            className: 'text-[var(--text)] border-l-[var(--text)]',
          }}
        >
          {item.label}
        </Link>
      ))}
    </>
  )
}
