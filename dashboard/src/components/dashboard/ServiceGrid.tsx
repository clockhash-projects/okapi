import type { ServiceRecord } from '@/types/health'
import { ServiceCard } from './ServiceCard'

interface ServiceGridProps {
  results: ServiceRecord[]
  onOpen: (record: ServiceRecord) => void
}

export function ServiceGrid({ results, onOpen }: ServiceGridProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
      {results.map((record) => (
        <ServiceCard key={record.service} record={record} onOpen={onOpen} />
      ))}
    </div>
  )
}
