// Restaurant detail page — shows restaurant info + menu

import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useRestaurant } from '@food-platform/hooks'
import { AppLayout } from '../../components/AppLayout'
import { Skeleton } from '@food-platform/ui'
import { Badge } from '@food-platform/ui'
import { formatEGP, formatDuration } from '@food-platform/utils'
import { ItemDetailModal } from '../../components/ItemDetailModal'

export function RestaurantDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data: restaurant, isLoading, error } = useRestaurant(id ?? '')

  const [selectedItemIndex, setSelectedItemIndex] = useState<{
    catIndex: number
    itemIndex: number
  } | null>(null)

  if (isLoading) {
    return (
      <AppLayout showHeader={false} showBottomNav={false}>
        <RestaurantDetailSkeleton />
      </AppLayout>
    )
  }

  if (error || !restaurant) {
    return (
      <AppLayout showHeader={false} showBottomNav={false}>
        <div className="min-h-screen flex items-center justify-center p-6">
          <div className="text-center">
            <span className="material-symbols-rounded text-error text-5xl mb-3">error</span>
            <p className="text-body text-text-secondary mb-4">
              {error?.message ?? 'المطعم مش موجود'}
            </p>
            <button
              onClick={() => navigate(-1)}
              className="text-primary font-semibold hover:underline"
            >
              رجوع
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  const selectedItem = selectedItemIndex
    ? restaurant.menu_categories[selectedItemIndex.catIndex]?.items[selectedItemIndex.itemIndex]
    : null

  return (
    <AppLayout showHeader={false} showBottomNav={false}>
      <div className="min-h-screen pb-24">
        {/* Cover image */}
        <div className="relative w-full h-48 md:h-64 bg-bg-tertiary">
          {restaurant.cover_url ? (
            <img
              src={restaurant.cover_url}
              alt={restaurant.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <span className="material-symbols-rounded text-text-tertiary text-6xl">
                restaurant
              </span>
            </div>
          )}

          {/* Gradient overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />

          {/* Back button */}
          <button
            onClick={() => navigate(-1)}
            className="absolute top-4 right-4 w-10 h-10 bg-white/90 rounded-full flex items-center justify-center shadow-md hover:bg-white transition-colors"
            aria-label="رجوع"
          >
            <span className="material-symbols-rounded text-text-primary">arrow_forward</span>
          </button>

          {/* Favorite + Share */}
          <div className="absolute top-4 left-4 flex gap-2">
            <button
              className="w-10 h-10 bg-white/90 rounded-full flex items-center justify-center shadow-md hover:bg-white transition-colors"
              aria-label="حفظ"
            >
              <span className="material-symbols-rounded text-text-primary">favorite</span>
            </button>
            <button
              className="w-10 h-10 bg-white/90 rounded-full flex items-center justify-center shadow-md hover:bg-white transition-colors"
              aria-label="مشاركة"
            >
              <span className="material-symbols-rounded text-text-primary">share</span>
            </button>
          </div>
        </div>

        {/* Restaurant info */}
        <div className="max-w-4xl mx-auto px-4 -mt-8 relative">
          <div className="bg-surface rounded-xl shadow-lg p-5">
            <h1 className="text-h2 font-bold text-text-primary">{restaurant.name}</h1>

            {/* Rating + cuisine */}
            <div className="flex items-center gap-3 mt-2 text-body-sm">
              <span className="flex items-center gap-1">
                <span className="text-warning">★</span>
                <span className="font-semibold text-text-primary">{restaurant.rating.toFixed(1)}</span>
                <span className="text-text-tertiary">({restaurant.rating_count} تقييم)</span>
              </span>
              <span className="text-text-tertiary">•</span>
              <span className="text-text-secondary">
                {restaurant.cuisine_types.join(' • ')}
              </span>
            </div>

            {/* ETA + distance + fee */}
            <div className="flex items-center gap-4 mt-3 text-caption text-text-secondary">
              <span className="flex items-center gap-1">
                <span className="material-symbols-rounded text-sm">schedule</span>
                {restaurant.eta_minutes_min}-{restaurant.eta_minutes_max} دقيقة
              </span>
              <span className="flex items-center gap-1">
                <span className="material-symbols-rounded text-sm">delivery_dining</span>
                {formatEGP(restaurant.delivery_fee)} توصيل
              </span>
              <span className="flex items-center gap-1">
                <span className="material-symbols-rounded text-sm">payments</span>
                {'💰'.repeat(restaurant.price_range)} للشخصين
              </span>
            </div>

            {/* Promo */}
            {restaurant.promo && (
              <div className="mt-3 bg-primary/5 border border-primary/20 rounded-lg p-3">
                <p className="text-body-sm text-primary font-semibold">
                  🎁 {restaurant.promo.title}
                </p>
                {restaurant.promo.description && (
                  <p className="text-caption text-text-secondary mt-0.5">
                    {restaurant.promo.description}
                  </p>
                )}
              </div>
            )}

            {/* Closed banner */}
            {!restaurant.is_open && (
              <div className="mt-3 bg-error/5 border border-error/20 rounded-lg p-3">
                <p className="text-body-sm text-error font-semibold">
                  المطعم مقفل دلوقتي
                </p>
              </div>
            )}
          </div>
        </div>

        {/* Menu */}
        <div className="max-w-4xl mx-auto px-4 mt-6">
          {/* Search menu */}
          <div className="sticky top-0 z-10 bg-bg-primary py-3">
            <div className="relative">
              <span className="material-symbols-rounded absolute right-3 top-1/2 -translate-y-1/2 text-text-tertiary text-xl">
                search
              </span>
              <input
                type="text"
                placeholder="ابحث في المنيو..."
                className="w-full h-12 pr-11 pl-4 rounded-lg border border-border bg-surface text-body placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
              />
            </div>
          </div>

          {/* Menu categories */}
          {restaurant.menu_categories.map((category, catIndex) => (
            <div key={category.id} className="mt-6">
              <h2 className="text-h3 font-bold text-text-primary mb-3 sticky top-16 bg-bg-primary py-2 z-10">
                {category.name}
              </h2>

              <div className="space-y-3">
                {category.items.map((item, itemIndex) => (
                  <MenuItemCard
                    key={item.id}
                    item={item}
                    onClick={() => setSelectedItemIndex({ catIndex, itemIndex })}
                  />
                ))}
              </div>
            </div>
          ))}

          {restaurant.menu_categories.length === 0 && (
            <div className="text-center py-12">
              <span className="material-symbols-rounded text-text-tertiary text-5xl mb-3">
                menu_book
              </span>
              <p className="text-body text-text-secondary">المنيو فاضي</p>
            </div>
          )}
        </div>
      </div>

      {/* Item detail modal */}
      {selectedItem && (
        <ItemDetailModal
          item={selectedItem}
          onClose={() => setSelectedItemIndex(null)}
        />
      )}
    </AppLayout>
  )
}

// ============ Menu Item Card ============

interface MenuItemCardProps {
  item: {
    id: string
    name: string
    description?: string
    price: string | number
    image_url?: string
    is_available: boolean
    prep_time_minutes: number
    rating?: number
    is_most_ordered?: boolean
  }
  onClick?: () => void
}

function MenuItemCard({ item, onClick }: MenuItemCardProps) {
  return (
    <div
      onClick={item.is_available ? onClick : undefined}
      className={`flex gap-3 p-3 bg-surface rounded-lg border border-border transition-all
        ${item.is_available
          ? 'cursor-pointer hover:shadow-md hover:border-primary/30'
          : 'opacity-50 cursor-not-allowed'
        }`}
    >
      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-start gap-2">
          <h3 className="text-body font-semibold text-text-primary">{item.name}</h3>
          {item.is_most_ordered && (
            <Badge variant="warning">🔥 الأكثر طلباً</Badge>
          )}
        </div>

        {item.description && (
          <p className="text-caption text-text-secondary mt-1 line-clamp-2">
            {item.description}
          </p>
        )}

        <div className="flex items-center gap-3 mt-2">
          <span className="text-body font-bold text-text-primary">
            {formatEGP(item.price)}
          </span>
          {item.rating && (
            <span className="text-caption text-text-tertiary flex items-center gap-0.5">
              <span className="text-warning">★</span>
              {item.rating.toFixed(1)}
            </span>
          )}
          <span className="text-caption text-text-tertiary">
            {formatDuration(item.prep_time_minutes)}
          </span>
        </div>

        {!item.is_available && (
          <span className="inline-block mt-2 text-caption text-error font-semibold">
            نفد
          </span>
        )}
      </div>

      {/* Image + add button */}
      <div className="relative w-24 h-24 flex-shrink-0">
        {item.image_url ? (
          <img
            src={item.image_url}
            alt={item.name}
            className="w-full h-full object-cover rounded-lg"
            loading="lazy"
          />
        ) : (
          <div className="w-full h-full bg-bg-tertiary rounded-lg flex items-center justify-center">
            <span className="material-symbols-rounded text-text-tertiary text-2xl">
              lunch_dining
            </span>
          </div>
        )}
        {item.is_available && (
          <button
            onClick={(e) => {
              e.stopPropagation()
              onClick?.()
            }}
            className="absolute -bottom-2 -left-2 w-9 h-9 bg-primary text-white rounded-full flex items-center justify-center shadow-md hover:bg-primary-dark transition-colors"
            aria-label="إضافة"
          >
            <span className="material-symbols-rounded text-xl">add</span>
          </button>
        )}
      </div>
    </div>
  )
}

// ============ Skeleton ============

function RestaurantDetailSkeleton() {
  return (
    <div className="min-h-screen pb-24">
      <Skeleton width="100%" height={200} rounded="md" />
      <div className="max-w-4xl mx-auto px-4 -mt-8 relative">
        <div className="bg-surface rounded-xl shadow-lg p-5">
          <Skeleton width="60%" height={28} />
          <div className="flex gap-3 mt-3">
            <Skeleton width={80} height={16} />
            <Skeleton width={60} height={16} />
          </div>
          <div className="flex gap-4 mt-4">
            <Skeleton width={100} height={14} />
            <Skeleton width={80} height={14} />
            <Skeleton width={60} height={14} />
          </div>
        </div>

        <div className="mt-6 space-y-4">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="flex gap-3 p-3 bg-surface rounded-lg border border-border">
              <div className="flex-1 space-y-2">
                <Skeleton width="50%" height={16} />
                <Skeleton width="80%" height={12} />
                <Skeleton width="30%" height={14} />
              </div>
              <Skeleton width={96} height={96} rounded="md" />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
