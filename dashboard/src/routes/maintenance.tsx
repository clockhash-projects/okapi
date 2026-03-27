import { createFileRoute, defer, Await } from '@tanstack/react-router'
import { okapi } from '@/lib/okapi'
import { StatusPill } from '@/components/shared/StatusPill'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBar } from '@/components/shared/FilterBar'
import { MaintenanceList } from '@/components/maintenance/MaintenanceList'
import { MaintenanceSkeleton } from '@/components/maintenance/MaintenanceSkeleton'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'
import { Pagination } from '@/components/shared/Pagination'
import { Suspense, useState } from 'react'
import type { MaintenanceResponse } from '@/types/maintenance'

interface MaintenanceSearch {
  limit?: number
  offset?: number
}

export const Route = createFileRoute('/maintenance')({
  validateSearch: (search: Record<string, unknown>): MaintenanceSearch => ({
    limit: Number(search.limit) || 20,
    offset: Number(search.offset) || 0,
  }),
  loaderDeps: ({ search: { limit, offset } }) => ({ limit, offset }),
  loader: ({ deps: { limit, offset } }) => ({
    maintenance: defer(okapi.maintenance({ limit, offset })),
  }),
  staticData: {
    title: 'Maintenance',
  },
  pendingMs: 0,
  pendingComponent: MaintenanceSkeleton,
  component: MaintenanceWrapper,
})

function MaintenanceWrapper() {
  const { maintenance } = Route.useLoaderData()

  return (
    <Suspense fallback={<MaintenanceSkeleton />}>
      <Await promise={maintenance}>
        {(data) => <MaintenanceView maintenance={data} />}
      </Await>
    </Suspense>
  )
}

const SECTION_TABS = ['Active', 'Upcoming', 'Completed'] as const
type SectionTab = (typeof SECTION_TABS)[number]

function MaintenanceView({ maintenance }: { maintenance: MaintenanceResponse }) {
  const allRecords = maintenance?.scheduled_maintenance || []
  const { limit = 20, offset = 0 } = Route.useSearch()
  const navigate = Route.useNavigate()
  const [tab, setTab] = useState<SectionTab>('Active')
  const [searchQ, setSearchQ] = useState('')

  // Each record has a list of maintenance windows; classify by window status
  const now = new Date()

  function classifyRecord(record: any): SectionTab | null {
    const windows = record.scheduled_maintenance || []
    if (windows.some((w: any) => w.status === 'in_progress')) return 'Active'
    if (windows.some((w: any) => w.status === 'scheduled' && new Date(w.starts_at) > now))
      return 'Upcoming'
    if (windows.every((w: any) => w.status === 'completed')) return 'Completed'
    return 'Upcoming' // fallback
  }

  const counts: Record<string, number> = {
    Active: allRecords.filter((r) => classifyRecord(r) === 'Active').length,
    Upcoming: allRecords.filter((r) => classifyRecord(r) === 'Upcoming').length,
    Completed: allRecords.filter((r) => classifyRecord(r) === 'Completed').length,
  }

  const filtered = allRecords.filter((r) => {
    const matchesTab = classifyRecord(r) === tab
    const matchesSearch = !searchQ || r.service.toLowerCase().includes(searchQ.toLowerCase())
    return matchesTab && matchesSearch
  })

  return (
    <div className="flex flex-col">
      <SectionLabel count={maintenance.count || 0}>Maintenance</SectionLabel>

      <FilterBar
        search={searchQ}
        onSearchChange={setSearchQ}
        status={tab}
        onStatusChange={setTab}
        statuses={SECTION_TABS}
        counts={counts}
        placeholder="SEARCH SERVICE..."
        label="Filter Maintenance"
      />

      {/* List */}
      {filtered.length > 0 ? (
        <div className="flex flex-col gap-4">
          <div className="border-t border-[var(--border)]">
            <Accordion multiple className="w-full">
              {filtered.map((record) => (
                <AccordionItem
                  key={record.service}
                  value={record.service}
                  className="border-b border-[var(--border)]"
                >
                  <AccordionTrigger className="hover:no-underline hover:bg-[var(--surface-2)] px-4 py-4 transition-colors">
                    <div className="flex items-center gap-4 w-full text-left">
                      <StatusPill status={record.status} />
                      <span className="font-mono text-sm text-[var(--text)]">{record.service}</span>
                      <div className="ml-auto flex items-center gap-6 pr-4">
                        <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase">
                          {(record.scheduled_maintenance || []).length} windows
                        </span>
                        <span className="font-mono text-[10px] text-[var(--text-muted)] uppercase">
                          {record.data_source}
                        </span>
                      </div>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent className="bg-[var(--surface)] p-6">
                    <div className="flex flex-col gap-3">
                      <span className="font-mono text-[10px] text-[var(--text-muted)] tracking-widest uppercase">
                        / Maintenance Windows
                      </span>
                      <MaintenanceList windows={record.scheduled_maintenance} />
                    </div>
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          </div>

          <Pagination
            total={maintenance.count || 0}
            limit={limit}
            offset={offset}
            onPageChange={(newOffset) =>
              navigate({ search: (prev) => ({ ...prev, offset: newOffset }) })
            }
          />
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-24 gap-4">
          <div className="w-2 h-2 rounded-full bg-[var(--text-muted)] opacity-50" />
          <span className="font-mono text-sm text-[var(--text-muted)]">
            No {tab.toLowerCase()} maintenance windows
          </span>
        </div>
      )}
    </div>
  )
}
