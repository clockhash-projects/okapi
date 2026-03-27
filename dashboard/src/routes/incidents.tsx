import { createFileRoute, defer, Await } from '@tanstack/react-router'
import { okapi } from '@/lib/okapi'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBar } from '@/components/shared/FilterBar'
import { IncidentAccordion } from '@/components/incidents/IncidentAccordion'
import { IncidentsSkeleton } from '@/components/incidents/IncidentsSkeleton'
import { Suspense, useState } from 'react'
import type { IncidentsResponse } from '@/types/incidents'

export const Route = createFileRoute('/incidents')({
  loader: () => ({
    incidents: defer(okapi.incidents()),
  }),
  staticData: {
    title: 'Incidents',
  },
  component: IncidentsWrapper,
})

function IncidentsWrapper() {
  const { incidents } = Route.useLoaderData()

  return (
    <Suspense fallback={<IncidentsSkeleton />}>
      <Await promise={incidents}>
        {(data) => <IncidentsView incidents={data} />}
      </Await>
    </Suspense>
  )
}

const STATUS_FILTERS = ['All', 'Investigating', 'Identified', 'Monitoring', 'Resolved'] as const
type StatusFilter = (typeof STATUS_FILTERS)[number]

function IncidentsView({ incidents }: { incidents: IncidentsResponse }) {
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('All')
  const [searchQ, setSearchQ] = useState('')

  const allAffected = (incidents?.incidents || []).filter(
    (r) => (r.incidents || []).length > 0
  )

  const counts: Record<string, number> = {
    All: allAffected.length,
    Investigating: allAffected.filter((s) =>
      (s.incidents || []).some((i) => i.status === 'investigating')
    ).length,
    Identified: allAffected.filter((s) =>
      (s.incidents || []).some((i) => i.status === 'identified')
    ).length,
    Monitoring: allAffected.filter((s) =>
      (s.incidents || []).some((i) => i.status === 'monitoring')
    ).length,
    Resolved: allAffected.filter((s) =>
      (s.incidents || []).some((i) => i.status === 'resolved')
    ).length,
  }

  // Client-side filter
  const filtered = allAffected.filter((r) => {
    const matchesStatus =
      statusFilter === 'All' ||
      (r.incidents || []).some(
        (i) => i.status === statusFilter.toLowerCase()
      )
    const matchesSearch =
      !searchQ || r.service.toLowerCase().includes(searchQ.toLowerCase())
    return matchesStatus && matchesSearch
  })

  // Sort: active first, resolved last
  const statusPriority: Record<string, number> = {
    investigating: 0,
    identified: 1,
    monitoring: 2,
    resolved: 3,
  }
  const sorted = [...filtered].sort((a, b) => {
    const aPriority = Math.min(
      ...(a.incidents || []).map((i) => statusPriority[i.status] ?? 4)
    )
    const bPriority = Math.min(
      ...(b.incidents || []).map((i) => statusPriority[i.status] ?? 4)
    )
    return aPriority - bPriority
  })

  return (
    <div className="flex flex-col">
      <SectionLabel count={allAffected.length}>Incidents</SectionLabel>
      
      <FilterBar
        search={searchQ}
        onSearchChange={setSearchQ}
        status={statusFilter}
        onStatusChange={setStatusFilter}
        statuses={STATUS_FILTERS}
        counts={counts}
        placeholder="SEARCH SERVICE..."
        label="Filter Incidents"
      />

      {/* Incident list */}
      <section>
        {sorted.length > 0 ? (
          <IncidentAccordion records={sorted} />
        ) : (
          <div className="font-mono text-sm text-[var(--text-muted)] py-12 text-center">
            {allAffected.length === 0
              ? 'No active incidents reported.'
              : 'No incidents match your filter.'}
          </div>
        )}
      </section>
    </div>
  )
}
