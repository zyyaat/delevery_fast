// Restaurant-related hooks

import { useQuery } from '@tanstack/react-query'
import { getApiClient } from '@food-platform/api-client'
import type { Restaurant, NearbyRestaurantsResponse, RestaurantDetail } from '@food-platform/types'

const QUERY_KEYS = {
  nearby: (lat: number, lng: number, radius: number) => ['restaurants', 'nearby', lat, lng, radius] as const,
  detail: (id: string) => ['restaurants', 'detail', id] as const,
  search: (q: string, lat: number, lng: number) => ['restaurants', 'search', q, lat, lng] as const,
}

export function useNearbyRestaurants(lat: number, lng: number, radius = 3) {
  return useQuery<NearbyRestaurantsResponse>({
    queryKey: QUERY_KEYS.nearby(lat, lng, radius),
    queryFn: () =>
      getApiClient().get(`/restaurants/nearby?lat=${lat}&lng=${lng}&radius=${radius}`),
    enabled: lat !== 0 && lng !== 0,
    staleTime: 60_000, // 1 min
  })
}

export function useRestaurant(id: string) {
  return useQuery<RestaurantDetail>({
    queryKey: QUERY_KEYS.detail(id),
    queryFn: () => getApiClient().get(`/restaurants/${id}`),
    enabled: !!id,
    staleTime: 5 * 60_000, // 5 min
  })
}

export function useSearchRestaurants(query: string, lat: number, lng: number) {
  return useQuery<NearbyRestaurantsResponse>({
    queryKey: QUERY_KEYS.search(query, lat, lng),
    queryFn: () =>
      getApiClient().get(`/restaurants/search?q=${encodeURIComponent(query)}&lat=${lat}&lng=${lng}`),
    enabled: query.length >= 2,
    staleTime: 30_000,
  })
}
