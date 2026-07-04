import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import { initApiClient, initWSClient } from '@food-platform/api-client'
import { zustandTokenStorage, useAuthStore } from '@food-platform/auth'
import '@food-platform/theme/globals.css'

// Material Symbols
const link = document.createElement('link')
link.href =
  'https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200'
link.rel = 'stylesheet'
document.head.appendChild(link)

// Initialize API client
const apiURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const wsURL = import.meta.env.VITE_WS_URL || 'ws://localhost:8081'

initApiClient({
  baseURL: apiURL,
  tokenStorage: zustandTokenStorage,
  refreshTokenUrl: '/api/v1/auth/refresh',
})

// Initialize WebSocket if authenticated
const token = useAuthStore.getState().accessToken
if (token) {
  initWSClient({
    url: wsURL,
    token,
  })
}

// React Query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: (failureCount, error) => {
        // Don't retry on 4xx errors
        if (error instanceof Error && 'statusCode' in error) {
          const status = (error as { statusCode: number }).statusCode
          if (status >= 400 && status < 500) return false
        }
        return failureCount < 3
      },
      refetchOnWindowFocus: false,
    },
  },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
