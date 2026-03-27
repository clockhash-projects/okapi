import { Component, type ErrorInfo, type ReactNode } from 'react'
import { AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface Props {
  children?: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false
  }

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo)
  }

  public render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback
      }

      return (
        <div className="flex flex-col items-center justify-center min-h-[400px] p-8 text-center bg-zinc-50 dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800">
          <div className="p-3 mb-4 bg-red-100 dark:bg-red-900/30 rounded-full">
            <AlertTriangle className="w-8 h-8 text-red-600 dark:text-red-400" />
          </div>
          <h2 className="mb-2 text-xl font-semibold text-zinc-900 dark:text-zinc-100">
            Something went wrong
          </h2>
          <p className="mb-6 text-zinc-600 dark:text-zinc-400 max-w-md">
            We encountered an unexpected error while rendering this component.
            {this.state.error && (
              <span className="block mt-2 text-sm font-mono bg-zinc-100 dark:bg-zinc-950 p-2 rounded">
                {this.state.error.message}
              </span>
            )}
          </p>
          <Button 
            onClick={() => window.location.reload()}
            variant="default"
          >
            Reload Page
          </Button>
        </div>
      )
    }

    return this.props.children
  }
}
