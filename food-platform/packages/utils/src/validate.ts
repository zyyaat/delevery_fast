// Validation utilities

/**
 * Validate Egyptian phone number
 * @example isPhoneValid("01012345678") → true
 */
export function isPhoneValid(phone: string): boolean {
  const cleaned = phone.replace(/\D/g, '')
  return cleaned.length === 11 && cleaned.startsWith('01') && /^01[0-2,5][0-9]{8}$/.test(cleaned)
}

/**
 * Validate email
 */
export function isEmailValid(email: string): boolean {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return re.test(email)
}

/**
 * Validate UUID
 */
export function isUUID(value: string): boolean {
  const re = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
  return re.test(value)
}

/**
 * Validate Egyptian National ID (14 digits)
 */
export function isNationalIdValid(id: string): boolean {
  const cleaned = id.replace(/\D/g, '')
  return cleaned.length === 14
}

/**
 * Validate OTP code
 */
export function isOtpValid(code: string): boolean {
  return /^\d{6}$/.test(code)
}

/**
 * Validate price (positive number with max 2 decimals)
 */
export function isPriceValid(price: string | number): boolean {
  const value = typeof price === 'string' ? parseFloat(price) : price
  return !isNaN(value) && value >= 0 && value < 1000000
}

/**
 * Validate coordinates
 */
export function areCoordinatesValid(lat: number, lng: number): boolean {
  return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}
