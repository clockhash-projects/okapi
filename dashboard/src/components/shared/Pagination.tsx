import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react'
import { clsx } from 'clsx'

interface PaginationProps {
  total: number
  limit: number
  offset: number
  onPageChange: (newOffset: number) => void
  className?: string
}

export function Pagination({ total, limit, offset, onPageChange, className }: PaginationProps) {
  const totalPages = Math.ceil(total / limit)
  const currentPage = Math.floor(offset / limit) + 1

  if (totalPages <= 1) return null

  const getPages = () => {
    const pages = []
    const maxVisible = 5
    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2))
    let end = Math.min(totalPages, start + maxVisible - 1)

    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1)
    }

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
    return pages
  }

  const handlePageClick = (page: number) => {
    onPageChange((page - 1) * limit)
  }

  return (
    <div className={clsx("flex flex-col sm:flex-row items-center justify-center gap-4 mt-8 pt-4 border-t border-[var(--border)]", className)}>
      <div className="flex items-center gap-1">
        <PageButton
          onClick={() => handlePageClick(1)}
          disabled={currentPage === 1}
          icon={<ChevronsLeft size={14} />}
          label="FIRST"
          hideLabelOnMobile
        />
        <PageButton
          onClick={() => handlePageClick(currentPage - 1)}
          disabled={currentPage === 1}
          icon={<ChevronLeft size={14} />}
          label="PREV"
        />
      </div>

      <div className="flex items-center gap-1 bg-[var(--surface-2)] p-1 rounded-[2px] border border-[var(--border)]">
        {getPages().map(page => (
          <button
            key={page}
            onClick={() => handlePageClick(page)}
            className={clsx(
              "min-w-[32px] h-8 px-2 flex flex-col items-center justify-center font-mono text-xs transition-all rounded-[1px] relative outline-none focus:outline-none focus:ring-0 leading-none",
              currentPage === page
                ? "bg-[var(--surface-3)] text-[var(--text)] font-bold border border-[var(--border-hover)]"
                : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-[var(--surface-3)]"
            )}
          >
            {page}
          </button>
        ))}
      </div>
      <div className="flex items-center gap-1">
        <PageButton
          onClick={() => handlePageClick(currentPage + 1)}
          disabled={currentPage === totalPages}
          icon={<ChevronRight size={14} />}
          label="NEXT"
          iconRight
        />
        <PageButton
          onClick={() => handlePageClick(totalPages)}
          disabled={currentPage === totalPages}
          icon={<ChevronsRight size={14} />}
          label="LAST"
          iconRight
          hideLabelOnMobile
        />
      </div>
    </div>
  )
}

interface PageButtonProps {
  onClick: () => void
  disabled: boolean
  icon: React.ReactNode
  label?: string
  iconRight?: boolean
  hideLabelOnMobile?: boolean
}

function PageButton({ onClick, disabled, icon, label, iconRight, hideLabelOnMobile }: PageButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={clsx(
        "h-8 px-3 flex items-center justify-center gap-2 border transition-all rounded-[2px] font-mono text-[10px] tracking-tighter uppercase",
        disabled
          ? "opacity-20 cursor-not-allowed border-[var(--border)]"
          : "bg-[var(--surface-2)] text-[var(--text)] border-[var(--border)] hover:bg-[var(--surface-3)] hover:border-[var(--border-hover)] active:translate-y-[1px]"
      )}
    >
      {!iconRight && icon}
      {label && (
        <span className={clsx(hideLabelOnMobile && "hidden md:inline")}>
          {label}
        </span>
      )}
      {iconRight && icon}
    </button>
  )
}
