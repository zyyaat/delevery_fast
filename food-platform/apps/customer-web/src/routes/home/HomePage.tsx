// Home page — main restaurant discovery feed

import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNearbyRestaurants } from '@food-platform/hooks'
import { AppLayout } from '../../components/AppLayout'
import { RestaurantCard } from '../../components/RestaurantCard'
import { RestaurantListSkeleton } from '../../components/Skeletons'
import type { CuisineType } from '@food-platform/types'

// Default Cairo coordinates (will be replaced with user's location)
const DEFAULT_LAT = 30.0444
const DEFAULT_LNG = 31.2357

const CUISINES: { type: CuisineType; label: string; icon: string }[] = [
  { type: 'egyptian', label: 'مصري', icon: '🍽️' },
  { type: 'italian', label: 'إيطالي', icon: '🍕' },
  { type: 'asian', label: 'آسيوي', icon: '🍜' },
  { type: 'fast_food', label: 'وجبات سريعة', icon: '🍔' },
  { type: 'healthy', label: 'صحي', icon: '🥗' },
  { type: 'desserts', label: 'حلويات', icon: '🍰' },
  { type: 'indian', label: 'هندي', icon: '🍛' },
  { type: 'lebanese', label: 'لبناني', icon: '🧆' },
  { type: 'breakfast', label: 'فطار', icon: '🍳' },
  { type: 'coffee', label: 'قهوة', icon: '☕' },
]

export function HomePage() {
  const navigate = useNavigate()
  const [lat] = useState(DEFAULT_LAT)
  const [lng] = useState(DEFAULT_LNG)

  const { data: restaurants, isLoading, error } = useNearbyRestaurants(lat, lng, 5)

  // Split restaurants by sections
  const trending = restaurants?.restaurants?.slice(0, 10) ?? []
  const under30 = restaurants?.restaurants?.filter(
    (r) => r.eta_minutes_max <= 30,
  ) ?? []
  const topRated = [...(restaurants?.restaurants ?? [])]
    .sort((a, b) => b.rating - a.rating)
    .slice(0, 10)

  const handleRestaurantClick = (id: string) => {
    navigate(`/restaurants/${id}`)
  }

  return (
    <AppLayout>
      <div className="max-w-7xl mx-auto px-4 py-4 space-y-8">
        {/* Welcome banner (first-time users) */}
        <WelcomeBanner />

        {/* Trending Near You */}
        <section>
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-h3 font-bold text-text-primary flex items-center gap-2">
              <span>🔥</span> رائج قريب منك
            </h2>
            <button className="text-caption text-primary font-semibold hover:underline">
              عرض الكل
            </button>
          </div>

          {isLoading ? (
            <RestaurantListSkeleton count={4} />
          ) : error ? (
            <ErrorMessage message="حصل خطأ، اسحب للتحديث" />
          ) : trending.length === 0 ? (
            <EmptyState message="مفيش مطاعم في منطقتك دلوقتي" />
          ) : (
            <div className="flex gap-4 overflow-x-auto hide-scrollbar snap-x-mandatory pb-2 -mx-4 px-4">
              {trending.map((restaurant) => (
                <div key={restaurant.id} className="snap-start">
                  <RestaurantCard
                    restaurant={restaurant}
                    distanceKm={restaurant.distance_km}
                    onClick={() => handleRestaurantClick(restaurant.id)}
                  />
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Under 30 minutes */}
        {under30.length > 0 && (
          <section>
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-h3 font-bold text-text-primary flex items-center gap-2">
                <span>⚡</span> وصل في 30 دقيقة
              </h2>
            </div>
            <div className="flex gap-4 overflow-x-auto hide-scrollbar snap-x-mandatory pb-2 -mx-4 px-4">
              {under30.map((restaurant) => (
                <div key={restaurant.id} className="snap-start">
                  <RestaurantCard
                    restaurant={restaurant}
                    distanceKm={restaurant.distance_km}
                    onClick={() => handleRestaurantClick(restaurant.id)}
                  />
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Cuisines grid */}
        <section>
          <h2 className="text-h3 font-bold text-text-primary mb-3 flex items-center gap-2">
            <span>🍽️</span> أنواع المطابخ
          </h2>
          <div className="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-3">
            {CUISINES.map((cuisine) => (
              <button
                key={cuisine.type}
                onClick={() => navigate(`/search?cuisine=${cuisine.type}`)}
                className="flex flex-col items-center gap-2 p-3 bg-surface rounded-lg border border-border hover:shadow-md hover:-translate-y-0.5 transition-all duration-normal"
              >
                <span className="text-2xl">{cuisine.icon}</span>
                <span className="text-caption font-medium text-text-primary">
                  {cuisine.label}
                </span>
              </button>
            ))}
          </div>
        </section>

        {/* Top rated */}
        {topRated.length > 0 && (
          <section>
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-h3 font-bold text-text-primary flex items-center gap-2">
                <span>🏆</span> الأعلى تقييماً
              </h2>
            </div>
            <div className="flex gap-4 overflow-x-auto hide-scrollbar snap-x-mandatory pb-2 -mx-4 px-4">
              {topRated.map((restaurant) => (
                <div key={restaurant.id} className="snap-start">
                  <RestaurantCard
                    restaurant={restaurant}
                    distanceKm={restaurant.distance_km}
                    onClick={() => handleRestaurantClick(restaurant.id)}
                  />
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </AppLayout>
  )
}

// ============ Sub-components ============

function WelcomeBanner() {
  return (
    <div className="bg-gradient-to-l from-primary to-primary-dark rounded-2xl p-5 text-white">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-h3 font-bold">🎁 خصم 50% على أول طلب!</h2>
          <p className="text-body-sm opacity-90 mt-1">استخدم الكود: WELCOME50</p>
        </div>
        <span className="material-symbols-rounded text-4xl opacity-50">celebration</span>
      </div>
    </div>
  )
}

function ErrorMessage({ message }: { message: string }) {
  return (
    <div className="bg-error/5 border border-error/20 rounded-lg p-4 text-center">
      <span className="material-symbols-rounded text-error text-3xl mb-2">error</span>
      <p className="text-body-sm text-error">{message}</p>
    </div>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="text-center py-8">
      <span className="material-symbols-rounded text-text-tertiary text-5xl mb-3">
        restaurant_menu
      </span>
      <p className="text-body text-text-secondary">{message}</p>
    </div>
  )
}
