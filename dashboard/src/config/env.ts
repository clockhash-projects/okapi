
/// <reference types="vite/client" />

export const env = {
  // Use /api relative URL by default since the Go backend now prefixes all API routes
  okapiBaseUrl: import.meta.env.VITE_OKAPI_BASE_URL || '/api',
  isDev: import.meta.env.DEV,
  environment: import.meta.env.VITE_ENVIRONMENT || 'production',
} as const
