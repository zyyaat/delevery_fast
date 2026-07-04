// Order-related hooks

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getApiClient } from '@food-platform/api-client'
import type { Order, CreateOrderRequest, CreateOrderResponse } from '@food-platform/types'

const QUERY_KEYS = {
  order: (id: string) => ['order', id] as const,
  activeOrders: ['orders', 'active'] as const,
  orderHistory: (limit: number, offset: number) => ['orders', 'history', limit, offset] as const,
}

export function useOrder(orderId: string) {
  return useQuery<Order>({
    queryKey: QUERY_KEYS.order(orderId),
    queryFn: () => getApiClient().get(`/orders/${orderId}`),
    enabled: !!orderId,
    staleTime: 30_000, // 30s
  })
}

export function useActiveOrders() {
  return useQuery<Order[]>({
    queryKey: QUERY_KEYS.activeOrders,
    queryFn: () => getApiClient().get('/customers/me/orders/active'),
    refetchInterval: 30_000, // Refresh every 30s
  })
}

export function useOrderHistory(limit = 20, offset = 0) {
  return useQuery<Order[]>({
    queryKey: QUERY_KEYS.orderHistory(limit, offset),
    queryFn: () => getApiClient().get(`/customers/me/orders?limit=${limit}&offset=${offset}`),
  })
}

export function useCreateOrder() {
  const queryClient = useQueryClient()

  return useMutation<CreateOrderResponse, Error, CreateOrderRequest>({
    mutationFn: (data) => getApiClient().postWithIdempotency('/orders', data),
    onSuccess: (order) => {
      queryClient.setQueryData(QUERY_KEYS.order(order.id), order)
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.activeOrders })
    },
  })
}

export function useCancelOrder() {
  const queryClient = useQueryClient()

  return useMutation<unknown, Error, { orderId: string; reason: string }>({
    mutationFn: ({ orderId, reason }) =>
      getApiClient().post(`/orders/${orderId}/cancel`, { reason }),
    onSuccess: (_, { orderId }) => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.order(orderId) })
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.activeOrders })
    },
  })
}
