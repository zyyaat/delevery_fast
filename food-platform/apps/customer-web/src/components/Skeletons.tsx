// Skeleton loaders for restaurant cards

import { Skeleton } from '@food-platform/ui'

export function RestaurantCardSkeleton() {
  return (
    <div className="w-72 flex-shrink-0 bg-surface rounded-lg border border-border overflow-hidden">
      <Skeleton width="100%" height={144} rounded="none" />
      <div className="p-3 space-y-2">
        <Skeleton width="60%" height={16} />
        <Skeleton width="40%" height={12} />
        <div className="flex gap-3">
          <Skeleton width={50} height={12} />
          <Skeleton width={60} height={12} />
        </div>
        <Skeleton width="80%" height={12} />
      </div>
    </div>
  )
}

export function RestaurantListSkeleton({ count = 4 }: { count?: number }) {
  return (
    <div className="flex gap-4 overflow-hidden">
      {Array.from({ length: count }).map((_, i) => (
        <RestaurantCardSkeleton key={i} />
      ))}
    </div>
  )
}

export function HorizontalRestaurantCardSkeleton() {
  return (
    <div className="flex gap-3 p-3 bg-surface rounded-lg border border-border">
      <Skeleton width={96} height={96} rounded="md" />
      <div className="flex-1 space-y-2">
        <Skeleton width="50%" height={16} />
        <Skeleton width="30%" height={12} />
        <div className="flex gap-3">
          <Skeleton width={40} height={12} />
          <Skeleton width={50} height={12} />
          <Skeleton width={40} height={12} />
        </div>
        <Skeleton width="60%" height={12} />
      </div>
    </div>
  )
}
