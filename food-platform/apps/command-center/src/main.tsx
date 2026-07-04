import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import { initApiClient } from '@food-platform/api-client'
import { zustandTokenStorage } from '@food-platform/auth'
import '@food-platform/theme/globals.css'

const apiURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

initApiClient({
  baseURL: apiURL,
  tokenStorage: zustandTokenStorage,
  refreshTokenUrl: '/api/v1/auth/refresh',
})

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
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
