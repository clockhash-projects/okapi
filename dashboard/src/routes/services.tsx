import { createFileRoute, defer, Await } from '@tanstack/react-router'
import { okapi } from '@/lib/okapi'
import { SectionLabel } from '@/components/shared/SectionLabel'
import { FilterBar } from '@/components/shared/FilterBar'
import { ServiceTable } from '@/components/services/ServiceTable'
import { ServicesSkeleton } from '@/components/services/ServicesSkeleton'
import { useState, Suspense } from 'react'
import { ServiceDrawer } from '@/components/dashboard/ServiceDrawer'
import { Pagination } from '@/components/shared/Pagination'
import { usePinnedServices } from '@/hooks/usePinnedServices'
import type { ServiceRecord, HealthAllResponse } from '@/types/health'

type ServicesSearch = {
  q?: string
  status?: string
  limit?: number
  offset?: number
}

export const Route = createFileRoute('/services')({
  validateSearch: (search: Record<string, unknown>): ServicesSearch => {
    return {
      q: (search.q as string) || undefined,
      status: (search.status as string) || 'All',
      limit: Number(search.limit) || 20,
      offset: Number(search.offset) || 0,
    }
  },
  loaderDeps: ({ search: { q, status, limit, offset } }) => ({ q, status, limit, offset }),
  loader: ({ deps: { q, status, limit, offset } }) => {
    return {
      servicesData: defer(
        okapi.healthAll({ q, status, limit, offset })
      ),
    }
  },
  staticData: {
    title: 'Services',
  },
  pendingMs: 0,
  pendingComponent: ServicesSkeleton,
  component: ServicesWrapper,
})

function ServicesWrapper() {
  const { servicesData } = Route.useLoaderData()

  return (
    <Suspense fallback={<ServicesSkeleton />}>
      <Await promise={servicesData}>
        {(healthData: HealthAllResponse) => <ServicesView health={healthData} />}
      </Await>
    </Suspense>
  )
}

function ServicesView({ health }: { health: HealthAllResponse }) {
  const [selectedService, setSelectedService] = useState<ServiceRecord | null>(null)
  const healthResults = health?.results || []
  const { q: searchQ, status: searchStatus, limit = 20, offset = 0 } = Route.useSearch()
  const navigate = Route.useNavigate()
  const { togglePin, isPinned } = usePinnedServices()

  const [q, setQ] = useState(searchQ || '')
  const [status, setStatus] = useState(searchStatus || 'All')

  const handleSearchChange = (val: string) => {
    setQ(val)
    navigate({ search: (prev) => ({ ...prev, q: val || undefined, offset: 0 }), replace: true })
  }

  const handleStatusChange = (val: string) => {
    setStatus(val)
    navigate({ search: (prev) => ({ ...prev, status: val, offset: 0 }), replace: true })
  }

  // Pinned services shown separately (use all results from current page or filter live from cache)
  const pinnedResults = healthResults.filter((r) => isPinned(r.service))
  const hasPins = pinnedResults.length > 0

  return (
    <div className="flex flex-col">
      {/* Pinned section */}
      {hasPins && (
        <section className="mb-8">
          <SectionLabel count={pinnedResults.length}>Pinned</SectionLabel>
          <ServiceTable
            records={pinnedResults}
            onOpen={setSelectedService}
            onTogglePin={togglePin}
            isPinned={isPinned}
          />
        </section>
      )}

      <SectionLabel count={health.count || 0}>Services</SectionLabel>
      <FilterBar
        search={q}
        onSearchChange={handleSearchChange}
        status={status}
        onStatusChange={handleStatusChange}
        statuses={['All', 'Operational', 'Degraded', 'Outage']}
        counts={health.statusCounts}
        placeholder="FILTER SERVICES..."
      />

      <div className="flex flex-col gap-4">
        <ServiceTable
          records={healthResults}
          onOpen={setSelectedService}
          onTogglePin={togglePin}
          isPinned={isPinned}
        />

        <Pagination
          total={health.count || 0}
          limit={limit}
          offset={offset}
          onPageChange={(newOffset) =>
            navigate({ search: (prev) => ({ ...prev, offset: newOffset }) })
          }
        />
      </div>

      <ServiceDrawer record={selectedService} onClose={() => setSelectedService(null)} />
    </div>
  )
}
