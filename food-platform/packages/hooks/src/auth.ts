// Auth-related hooks

import { useMutation } from '@tanstack/react-query'
import { getApiClient } from '@food-platform/api-client'
import { useAuthStore } from '@food-platform/auth'
import type { SendOtpRequest, SendOtpResponse, VerifyOtpRequest, AuthResponse } from '@food-platform/types'

export function useSendOtp() {
  return useMutation<SendOtpResponse, Error, SendOtpRequest>({
    mutationFn: (data) => getApiClient().post('/auth/otp/send', data),
  })
}

export function useVerifyOtp() {
  const setAuth = useAuthStore((s) => s.setAuth)

  return useMutation<AuthResponse, Error, VerifyOtpRequest>({
    mutationFn: (data) => getApiClient().post('/auth/otp/verify', data),
    onSuccess: (auth) => {
      setAuth(auth)
    },
  })
}

export function useLogout() {
  const logout = useAuthStore((s) => s.logout)

  return useMutation<void, Error, void>({
    mutationFn: async () => {
      try {
        await getApiClient().post('/auth/logout')
      } catch {
        // Ignore errors — we're logging out anyway
      }
    },
    onSuccess: () => {
      logout()
    },
  })
}
