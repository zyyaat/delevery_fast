// API Client — Axios-based HTTP client with auth interceptors

import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type AxiosError,
  type InternalAxiosRequestConfig,
} from 'axios'
import type {
  AuthResponse,
  ErrorResponse,
} from '@food-platform/types'
import { uuid } from '@food-platform/utils'

// Token storage interface (implemented by auth package)
export interface TokenStorage {
  getAccessToken(): string | null
  getRefreshToken(): string | null
  setAuth(auth: AuthResponse): void
  clear(): void
  onTokenRefreshed?(token: string): void
}

export interface ApiClientConfig {
  baseURL: string
  timeout?: number
  tokenStorage: TokenStorage
  refreshTokenUrl?: string
}

export class ApiError extends Error {
  constructor(
    public readonly code: string,
    message: string,
    public readonly statusCode: number,
    public readonly details?: Record<string, unknown>,
    public readonly requestId?: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export class ApiClient {
  private client: AxiosInstance
  private tokenStorage: TokenStorage
  private refreshTokenUrl?: string
  private isRefreshing = false
  private refreshPromise: Promise<string> | null = null

  constructor(config: ApiClientConfig) {
    this.tokenStorage = config.tokenStorage
    this.refreshTokenUrl = config.refreshTokenUrl ?? undefined

    this.client = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout ?? 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor — attach auth + request ID
    this.client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
      const token = this.tokenStorage.getAccessToken()
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      config.headers['X-Request-ID'] = uuid()
      return config
    })

    // Response interceptor — handle errors + token refresh
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response.data,
      async (error: AxiosError<ErrorResponse>) => {
        const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

        // Try to refresh on 401
        if (
          error.response?.status === 401 &&
          !originalRequest._retry &&
          this.refreshTokenUrl &&
          this.tokenStorage.getRefreshToken()
        ) {
          originalRequest._retry = true

          try {
            const newToken = await this.refreshToken()
            originalRequest.headers = {
              ...originalRequest.headers,
              Authorization: `Bearer ${newToken}`,
            }
            return this.client(originalRequest)
          } catch (refreshError) {
            this.tokenStorage.clear()
            throw this.normalizeError(refreshError as AxiosError<ErrorResponse>)
          }
        }

        return Promise.reject(this.normalizeError(error))
      },
    )
  }

  private async refreshToken(): Promise<string> {
    // If already refreshing, wait for the existing promise
    if (this.isRefreshing && this.refreshPromise) {
      return this.refreshPromise
    }

    this.isRefreshing = true
    this.refreshPromise = this.doRefresh()

    try {
      const newToken = await this.refreshPromise
      return newToken
    } finally {
      this.isRefreshing = false
      this.refreshPromise = null
    }
  }

  private async doRefresh(): Promise<string> {
    const refreshToken = this.tokenStorage.getRefreshToken()
    if (!refreshToken || !this.refreshTokenUrl) {
      throw new ApiError('AUTH_REFRESH_INVALID', 'No refresh token', 401)
    }

    const response = await axios.post(
      this.refreshTokenUrl,
      { refresh_token: refreshToken },
      { baseURL: this.client.defaults.baseURL },
    )

    const auth = response.data as AuthResponse
    this.tokenStorage.setAuth(auth)
    this.tokenStorage.onTokenRefreshed?.(auth.access_token)
    return auth.access_token
  }

  private normalizeError(error: AxiosError<ErrorResponse>): ApiError {
    if (error.response) {
      const { data, status } = error.response
      if (data?.error) {
        return new ApiError(
          data.error.code,
          data.error.message,
          status,
          data.error.details,
          data.error.request_id,
        )
      }
      return new ApiError('INTERNAL_ERROR', 'An error occurred', status)
    }

    if (error.request) {
      return new ApiError('NETWORK_ERROR', 'Network error — check your connection', 0)
    }

    return new ApiError('UNKNOWN_ERROR', error.message || 'Unknown error', 0)
  }

  // ============ Public Methods ============

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.client.get(url, config)
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return this.client.post(url, data, config)
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return this.client.put(url, data, config)
  }

  async patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return this.client.patch(url, data, config)
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.client.delete(url, config)
  }

  // With idempotency key for state-mutating POSTs
  async postWithIdempotency<T>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig,
  ): Promise<T> {
    const idempotencyKey = uuid()
    return this.client.post(url, data, {
      ...config,
      headers: { ...config?.headers, 'X-Idempotency-Key': idempotencyKey },
    })
  }
}

// ============ Singleton instance ============

let _instance: ApiClient | null = null

export function initApiClient(config: ApiClientConfig): ApiClient {
  _instance = new ApiClient(config)
  return _instance
}

export function getApiClient(): ApiClient {
  if (!_instance) {
    throw new Error('ApiClient not initialized. Call initApiClient() first.')
  }
  return _instance
}
