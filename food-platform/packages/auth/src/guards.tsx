// Route guards for React Router

import type { JSX } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import type { UserRole } from '@food-platform/types'
import { useAuthStore } from './store'

interface RequireAuthProps {
  children: JSX.Element
  roles?: UserRole[]
}

/**
 * Guard that requires authentication
 * Optionally checks for specific roles
 */
export function RequireAuth({ children, roles }: RequireAuthProps) {
  const { isAuthenticated, user } = useAuthStore()
  const location = useLocation()

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  if (roles && !roles.includes(user.role)) {
    return <Navigate to="/forbidden" replace />
  }

  return children
}

/**
 * Guard that redirects authenticated users away from auth pages
 */
export function RedirectIfAuthed({ children }: { children: JSX.Element }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return children
}
