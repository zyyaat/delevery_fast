// Auth store — Zustand with persistence

import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import type { AuthResponse, User, UserRole } from '@food-platform/types'
import type { TokenStorage } from '@food-platform/api-client'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean

  // Actions
  setAuth: (auth: AuthResponse) => void
  setUser: (user: User) => void
  logout: () => void
  hasRole: (...roles: UserRole[]) => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      setAuth: (auth: AuthResponse) =>
        set({
          user: auth.user,
          accessToken: auth.access_token,
          refreshToken: auth.refresh_token,
          isAuthenticated: true,
        }),

      setUser: (user: User) => set({ user }),

      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        }),

      hasRole: (...roles: UserRole[]) => {
        const user = get().user
        return user ? roles.includes(user.role) : false
      },
    }),
    {
      name: 'food-platform-auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
)

// ============ Token Storage Adapter ============
// Bridges Zustand store with the ApiClient's TokenStorage interface

export const zustandTokenStorage: TokenStorage = {
  getAccessToken() {
    return useAuthStore.getState().accessToken
  },
  getRefreshToken() {
    return useAuthStore.getState().refreshToken
  },
  setAuth(auth: AuthResponse) {
    useAuthStore.getState().setAuth(auth)
  },
  clear() {
    useAuthStore.getState().logout()
  },
}
