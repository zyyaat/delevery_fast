// Crypto utilities

/**
 * Generate a random string (for IDs, nonces, etc.)
 */
export function randomString(length: number = 32): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  const cryptoObj = typeof crypto !== 'undefined' ? crypto : (undefined as unknown as Crypto)
  if (cryptoObj) {
    const values = new Uint32Array(length)
    cryptoObj.getRandomValues(values)
    for (let i = 0; i < length; i++) {
      result += chars[values[i]! % chars.length]
    }
  } else {
    // Fallback (not cryptographically secure)
    for (let i = 0; i < length; i++) {
      result += chars[Math.floor(Math.random() * chars.length)]!
    }
  }
  return result
}

/**
 * Generate a UUID v4
 */
export function uuid(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  // Fallback
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

/**
 * Hash a string (for fingerprinting, NOT for passwords)
 */
export async function hashString(input: string): Promise<string> {
  if (typeof crypto !== 'undefined' && crypto.subtle) {
    const encoder = new TextEncoder()
    const data = encoder.encode(input)
    const hashBuffer = await crypto.subtle.digest('SHA-256', data)
    return Array.from(new Uint8Array(hashBuffer))
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('')
  }
  // Simple fallback hash (NOT secure)
  let hash = 0
  for (let i = 0; i < input.length; i++) {
    const char = input.charCodeAt(i)
    hash = (hash << 5) - hash + char
    hash = hash & hash
  }
  return hash.toString(16)
}

/**
 * Generate device fingerprint (basic)
 */
export function getDeviceFingerprint(): string {
  const components = [
    navigator.userAgent,
    navigator.language,
    navigator.platform,
    screen.width + 'x' + screen.height,
    screen.colorDepth,
    new Date().getTimezoneOffset(),
  ]
  return components.join('|')
}
