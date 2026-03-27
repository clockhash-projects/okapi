import { useState, useEffect } from 'react'
import { okapi } from '@/lib/okapi'

export function useApiHealth() {
  const [reachable, setReachable] = useState<boolean | null>(null)
  const [latency, setLatency] = useState<number | null>(null)
  const [uptime, setUptime] = useState<string | null>(null)

  useEffect(() => {
    const checkHealth = async () => {
      const start = performance.now()
      try {
        const health = await okapi.selfHealth()
        setLatency(Math.round(performance.now() - start))
        setReachable(true)
        setUptime(health.uptime || null)
      } catch {
        setReachable(false)
        setLatency(null)
        setUptime(null)
      }
    }

    checkHealth()
    const interval = setInterval(checkHealth, 30000)
    return () => clearInterval(interval)
  }, [])

  return { reachable, latency, uptime }
}
