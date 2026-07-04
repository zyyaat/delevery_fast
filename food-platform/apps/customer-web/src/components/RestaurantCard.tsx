// Restaurant card component for home feed and search results

import { type Restaurant } from '@food-platform/types'
import { formatEGP, formatDistance, formatDuration } from '@food-platform/utils'
import { Badge } from '@food-platform/ui'
import { cn } from '@food-platform/ui'

interface RestaurantCardProps {
  restaurant: Restaurant
  distanceKm?: number
  onClick?: () => void
  variant?: 'horizontal' | 'vertical'
}

export function RestaurantCard({
  restaurant,
  distanceKm,
  onClick,
  variant = 'vertical',
}: RestaurantCardProps) {
  const isOpen = restaurant.is_open

  if (variant === 'horizontal') {
    return (
      <div
        onClick={onClick}
        className={cn(
          'flex gap-3 p-3 bg-surface rounded-lg border border-border cursor-pointer',
          'hover:shadow-md transition-all duration-normal',
          !isOpen && 'opacity-60',
        )}
      >
        {/* Image */}
        <div className="w-24 h-24 rounded-lg overflow-hidden bg-bg-tertiary flex-shrink-0">
          {restaurant.cover_url ? (
            <img
              src={restaurant.cover_url}
              alt={restaurant.name}
              className="w-full h-full object-cover"
              loading="lazy"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <span className="material-symbols-rounded text-text-tertiary text-3xl">restaurant</span>
            </div>
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <h3 className="text-body font-bold text-text-primary truncate">{restaurant.name}</h3>
            {restaurant.promo && (
              <Badge variant="primary">{restaurant.promo.title}</Badge>
            )}
          </div>
          <p className="text-caption text-text-secondary mt-0.5">
            {restaurant.cuisine_types.join(' • ')}
          </p>
          <div className="flex items-center gap-3 mt-2 text-caption text-text-tertiary">
            <span className="flex items-center gap-0.5">
              <span className="text-warning">★</span>
              {restaurant.rating.toFixed(1)}
              <span className="text-text-tertiary">({restaurant.rating_count})</span>
            </span>
            {distanceKm && <span>{formatDistance(distanceKm)}</span>}
            <span>{restaurant.eta_minutes_min}-{restaurant.eta_minutes_max} د</span>
          </div>
          <div className="flex items-center gap-2 mt-1">
            <span className="text-caption text-text-tertiary">
              {formatEGP(restaurant.delivery_fee)} توصيل
            </span>
            {!isOpen && <Badge variant="error">مقفل</Badge>}
          </div>
        </div>
      </div>
    )
  }

  // Vertical card (default for horizontal scroll)
  return (
    <div
      onClick={onClick}
      className={cn(
        'w-72 flex-shrink-0 bg-surface rounded-lg border border-border overflow-hidden cursor-pointer',
        'hover:shadow-lg hover:-translate-y-0.5 transition-all duration-normal',
        !isOpen && 'opacity-60',
      )}
    >
      {/* Image */}
      <div className="w-full h-36 bg-bg-tertiary relative">
        {restaurant.cover_url ? (
          <img
            src={restaurant.cover_url}
            alt={restaurant.name}
            className="w-full h-full object-cover"
            loading="lazy"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <span className="material-symbols-rounded text-text-tertiary text-4xl">restaurant</span>
          </div>
        )}
        {restaurant.promo && (
          <div className="absolute top-2 right-2">
            <Badge variant="primary">🎁 {restaurant.promo.title}</Badge>
          </div>
        )}
        {!isOpen && (
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
            <span className="text-white font-bold text-body">مقفل</span>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="p-3">
        <h3 className="text-body font-bold text-text-primary truncate">{restaurant.name}</h3>
        <p className="text-caption text-text-secondary mt-0.5 truncate">
          {restaurant.cuisine_types.join(' • ')}
        </p>
        <div className="flex items-center gap-3 mt-2 text-caption">
          <span className="flex items-center gap-0.5 text-text-primary">
            <span className="text-warning">★</span>
            <span className="font-semibold">{restaurant.rating.toFixed(1)}</span>
            <span className="text-text-tertiary">({restaurant.rating_count})</span>
          </span>
          <span className="text-text-tertiary">•</span>
          <span className="text-text-secondary">
            {restaurant.eta_minutes_min}-{restaurant.eta_minutes_max} دقيقة
          </span>
        </div>
        <div className="flex items-center justify-between mt-2">
          <span className="text-caption text-text-tertiary">
            {formatEGP(restaurant.delivery_fee)} توصيل
          </span>
          {distanceKm && (
            <span className="text-caption text-text-tertiary">{formatDistance(distanceKm)}</span>
          )}
        </div>
      </div>
    </div>
  )
}
