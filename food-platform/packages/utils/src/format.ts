// Formatting utilities

import type { Money } from '@food-platform/types'

/**
 * Format money in EGP
 * @example formatEGP("245.50") → "EGP 245.50"
 */
export function formatEGP(amount: Money | number): string {
  const value = typeof amount === 'string' ? parseFloat(amount) : amount
  return `EGP ${value.toFixed(2)}`
}

/**
 * Format distance in km
 * @example formatDistance(1.2) → "1.2 km"
 * @example formatDistance(0.8) → "800 m"
 */
export function formatDistance(km: number): string {
  if (km < 1) {
    return `${Math.round(km * 1000)} m`
  }
  return `${km.toFixed(1)} km`
}

/**
 * Format duration in minutes
 * @example formatDuration(35) → "35 دقيقة"
 * @example formatDuration(90) → "1 ساعة 30 دقيقة"
 */
export function formatDuration(minutes: number): string {
  if (minutes < 60) {
    return `${minutes} دقيقة`
  }
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  if (mins === 0) {
    return `${hours} ساعة`
  }
  return `${hours} ساعة ${mins} دقيقة`
}

/**
 * Format date in Arabic
 * @example formatDate("2026-07-04T14:32:15Z") → "4 يوليو 2026"
 */
export function formatDate(isoDate: string): string {
  const formatter = new Intl.DateTimeFormat('ar-EG', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
  return formatter.format(new Date(isoDate))
}

/**
 * Format time in Arabic
 * @example formatTime("2026-07-04T14:32:15Z") → "2:32 م"
 */
export function formatTime(isoDate: string): string {
  const formatter = new Intl.DateTimeFormat('ar-EG', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  })
  return formatter.format(new Date(isoDate))
}

/**
 * Format date and time together
 */
export function formatDateTime(isoDate: string): string {
  return `${formatDate(isoDate)} - ${formatTime(isoDate)}`
}

/**
 * Format relative time (time ago)
 * @example formatRelativeTime(past) → "منذ 5 دقائق"
 */
export function formatRelativeTime(isoDate: string): string {
  const now = Date.now()
  const then = new Date(isoDate).getTime()
  const diffSeconds = Math.floor((now - then) / 1000)

  if (diffSeconds < 60) return 'الآن'
  if (diffSeconds < 3600) return `منذ ${Math.floor(diffSeconds / 60)} دقيقة`
  if (diffSeconds < 86400) return `منذ ${Math.floor(diffSeconds / 3600)} ساعة`
  if (diffSeconds < 604800) return `منذ ${Math.floor(diffSeconds / 86400)} يوم`
  return formatDate(isoDate)
}

/**
 * Format countdown timer
 * @example formatCountdown(75) → "01:15"
 */
export function formatCountdown(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

/**
 * Format phone number (Egyptian)
 * @example formatPhone("01012345678") → "010 1234 5678"
 */
export function formatPhone(phone: string): string {
  const cleaned = phone.replace(/\D/g, '')
  if (cleaned.length === 11 && cleaned.startsWith('01')) {
    return `${cleaned.slice(0, 3)} ${cleaned.slice(3, 7)} ${cleaned.slice(7)}`
  }
  return phone
}

/**
 * Truncate text with ellipsis
 */
export function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength - 1) + '…'
}
