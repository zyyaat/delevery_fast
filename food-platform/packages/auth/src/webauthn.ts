// WebAuthn helpers for biometric authentication (employee portal)

import type { UserRole } from '@food-platform/types'

export interface WebAuthnRegistrationOptions {
  challenge: string
  rp: {
    name: string
    id: string
  }
  user: {
    id: string
    name: string
    displayName: string
  }
  pubKeyCredParams: {
    type: 'public-key'
    alg: number
  }[]
  authenticatorSelection: {
    authenticatorAttachment?: 'platform' | 'cross-platform'
    userVerification: 'required' | 'preferred' | 'discouraged'
  }
  timeout: number
}

export interface WebAuthnCredential {
  id: string
  rawId: string
  type: 'public-key'
  response: {
    attestationObject: string
    clientDataJSON: string
  }
}

export interface WebAuthnAssertion {
  id: string
  rawId: string
  type: 'public-key'
  response: {
    authenticatorData: string
    clientDataJSON: string
    signature: string
    userHandle?: string
  }
}

/**
 * Check if WebAuthn is available in this browser
 */
export function isWebAuthnAvailable(): boolean {
  return typeof window !== 'undefined' && 'PublicKeyCredential' in window
}

/**
 * Check if platform authenticator (Touch ID, Face ID, Windows Hello) is available
 */
export async function isPlatformAuthenticatorAvailable(): Promise<boolean> {
  if (!isWebAuthnAvailable()) return false
  try {
    return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable()
  } catch {
    return false
  }
}

/**
 * Base64URL encode
 */
export function base64urlEncode(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let str = ''
  for (const byte of bytes) {
    str += String.fromCharCode(byte)
  }
  return btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

/**
 * Base64URL decode
 */
export function base64urlDecode(str: string): ArrayBuffer {
  const decoded = atob(str.replace(/-/g, '+').replace(/_/g, '/'))
  const bytes = new Uint8Array(decoded.length)
  for (let i = 0; i < decoded.length; i++) {
    bytes[i] = decoded.charCodeAt(i)
  }
  return bytes.buffer
}

/**
 * Begin WebAuthn registration
 */
export async function beginRegistration(
  options: WebAuthnRegistrationOptions,
): Promise<WebAuthnCredential> {
  const publicKey: PublicKeyCredentialCreationOptions = {
    challenge: base64urlDecode(options.challenge),
    rp: options.rp,
    user: {
      id: base64urlDecode(options.user.id),
      name: options.user.name,
      displayName: options.user.displayName,
    },
    pubKeyCredParams: options.pubKeyCredParams,
    authenticatorSelection: {
      authenticatorAttachment: options.authenticatorSelection.authenticatorAttachment,
      userVerification: options.authenticatorSelection.userVerification,
      residentKey: 'preferred',
      requireResidentKey: false,
    },
    timeout: options.timeout,
    attestation: 'none',
  }

  const credential = (await navigator.credentials.create({ publicKey })) as PublicKeyCredential

  if (!credential) {
    throw new Error('Registration failed: no credential returned')
  }

  const response = credential.response as AuthenticatorAttestationResponse

  return {
    id: credential.id,
    rawId: base64urlEncode(credential.rawId),
    type: 'public-key',
    response: {
      attestationObject: base64urlEncode(response.attestationObject),
      clientDataJSON: base64urlEncode(response.clientDataJSON),
    },
  }
}

/**
 * Verify WebAuthn (assertion) for sensitive action
 */
export async function verifyWebAuthn(
  options: {
    challenge: string
    allowCredentials?: { id: string; type: 'public-key' }[]
  },
  action: string,
  context?: Record<string, unknown>,
): Promise<{ assertion: WebAuthnAssertion; action: string; context?: Record<string, unknown> }> {
  const publicKey: PublicKeyCredentialRequestOptions = {
    challenge: base64urlDecode(options.challenge),
    allowCredentials: options.allowCredentials?.map((c) => ({
      id: base64urlDecode(c.id),
      type: 'public-key',
    })),
    userVerification: 'required',
    timeout: 60000,
  }

  const assertion = (await navigator.credentials.get({ publicKey })) as PublicKeyCredential

  if (!assertion) {
    throw new Error('Verification failed: no assertion returned')
  }

  const response = assertion.response as AuthenticatorAssertionResponse

  return {
    assertion: {
      id: assertion.id,
      rawId: base64urlEncode(assertion.rawId),
      type: 'public-key',
      response: {
        authenticatorData: base64urlEncode(response.authenticatorData),
        clientDataJSON: base64urlEncode(response.clientDataJSON),
        signature: base64urlEncode(response.signature),
        userHandle: response.userHandle ? base64urlEncode(response.userHandle) : undefined,
      },
    },
    action,
    context,
  }
}

/**
 * Role-based access check
 */
export function canAccess(requiredRole: UserRole, userRole: UserRole): boolean {
  const roleHierarchy: Record<UserRole, number> = {
    read_only_analyst: 0,
    customer: 1,
    driver: 1,
    restaurant: 1,
    support_l1: 2,
    field_supervisor: 3,
    support_l2: 3,
    hr: 4,
    finance: 4,
    ops_manager: 5,
    super_admin: 10,
  }

  return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
