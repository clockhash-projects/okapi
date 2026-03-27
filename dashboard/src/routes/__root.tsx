import { createRootRoute, Outlet, useRouterState } from "@tanstack/react-router"
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools"
import { TanStackDevtools } from "@tanstack/react-devtools"
import { Sidebar } from "@/components/layout/Sidebar"
import { Header } from "@/components/layout/Header"
import { ErrorBoundary } from "@/components/shared/ErrorBoundary"
import { NotFound } from "@/components/shared/NotFound"
import { useMatches } from "@tanstack/react-router"

export const Route = createRootRoute({
  component: RootComponent,
  notFoundComponent: NotFound,
})

interface RouteData {
  title?: string
}

interface LoaderData {
  health?: { results: { fetched_at: string }[] }
  incidents?: { incidents: { fetched_at: string }[] }
  history?: { history: { fetched_at: string }[] }
}

function RouterProgressBar() {
  const isTransitioning = useRouterState({ select: (s) => s.isTransitioning })
  if (!isTransitioning) return null
  return <div className="nav-progress-bar" />
}

function RootComponent() {
  const matches = useMatches()
  const lastMatch = matches[matches.length - 1]

  const staticData = (lastMatch?.staticData || {}) as RouteData
  const loaderData = (lastMatch?.loaderData || {}) as LoaderData

  const title = staticData.title || 'Okapi'
  const fetchedAt = loaderData.health?.results?.[0]?.fetched_at ||
    loaderData.incidents?.incidents?.[0]?.fetched_at ||
    loaderData.history?.history?.[0]?.fetched_at

  return (
    <div className="bg-[var(--bg)] text-[var(--text)] min-h-screen">
      <RouterProgressBar />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 ml-0 lg:ml-60 min-h-screen flex flex-col">
          <Header title={title} lastFetched={fetchedAt} />
          <div className="p-4 sm:p-6 lg:p-8 w-full max-w-[1600px] mx-auto">
            <ErrorBoundary>
              <Outlet />
            </ErrorBoundary>
          </div>
        </main>
      </div>
      {process.env.NODE_ENV !== 'production' && (
        <TanStackDevtools
          config={{
            position: "bottom-right",
          }}
          plugins={[
            {
              name: "Tanstack Router",
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
      )}
    </div>
  )
}
