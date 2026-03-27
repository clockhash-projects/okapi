import { createRouter as createTanStackRouter } from "@tanstack/react-router"
import { routeTree } from "./routeTree.gen"

export function getRouter() {
  const router = createTanStackRouter({
    routeTree,

    scrollRestoration: true,

    // Preload routes when the user hovers a link 
    defaultPreload: "intent",
    defaultPreloadDelay: 100,        // Start preloading after 100ms hover (avoid wasted preloads)
    defaultPreloadStaleTime: 0,      // Always re-fetch stale data on navigate

    // Pending component timing (these were the bugs)
    // Official docs: defaultPendingMs defaults to 1000ms — that's why skeletons never appeared!
    defaultPendingMs: 0,             // Show pendingComponent immediately on navigate
    defaultPendingMinMs: 300,        // Keep it for at least 300ms to avoid flash
  })

  return router
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
