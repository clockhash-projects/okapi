import { useState, useCallback } from 'react'

const STORAGE_KEY = 'okapi:pinned'

function load(): Set<string> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return new Set(JSON.parse(raw) as string[])
  } catch {}
  return new Set()
}

function save(pins: Set<string>) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify([...pins]))
}

export function usePinnedServices() {
  const [pinned, setPinned] = useState<Set<string>>(load)

  const togglePin = useCallback((id: string) => {
    setPinned((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      save(next)
      return next
    })
  }, [])

  const isPinned = useCallback((id: string) => pinned.has(id), [pinned])

  return { pinned, togglePin, isPinned }
}
