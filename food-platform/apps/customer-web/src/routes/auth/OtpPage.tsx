// OTP verification page — customer enters 6-digit code

import { useState, useRef, useEffect } from 'react'
import { useNavigate, useLocation, Link } from 'react-router-dom'
import { useVerifyOtp } from '@food-platform/hooks'
import { Button } from '@food-platform/ui'
import { useCountdown } from '@food-platform/hooks'

export function OtpPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const verifyOtp = useVerifyOtp()

  const state = (location.state ?? {}) as {
    requestId?: string
    phone?: string
    expiresIn?: number
  }

  // Redirect to login if no request ID
  useEffect(() => {
    if (!state.requestId) {
      navigate('/login', { replace: true })
    }
  }, [state.requestId, navigate])

  const [digits, setDigits] = useState<string[]>(Array(6).fill(''))
  const [error, setError] = useState('')
  const inputsRef = useRef<(HTMLInputElement | null)[]>([])

  const { seconds, formatted, isRunning } = useCountdown(
    state.expiresIn ?? 120,
    () => {
      // Auto-submit when all digits entered
    },
  )

  // Auto-focus first input on mount
  useEffect(() => {
    inputsRef.current[0]?.focus()
  }, [])

  const handleChange = (index: number, value: string) => {
    // Only allow digits
    const digit = value.replace(/\D/g, '').slice(-1)

    const newDigits = [...digits]
    newDigits[index] = digit
    setDigits(newDigits)
    setError('')

    // Auto-advance to next input
    if (digit && index < 5) {
      inputsRef.current[index + 1]?.focus()
    }

    // Auto-submit when all 6 digits entered
    if (digit && index === 5) {
      const code = newDigits.join('')
      if (code.length === 6) {
        handleVerify(code)
      }
    }
  }

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace' && !digits[index] && index > 0) {
      // Move to previous input on backspace
      inputsRef.current[index - 1]?.focus()
    }
  }

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault()
    const pasted = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6)
    if (pasted.length > 0) {
      const newDigits = Array(6).fill('')
      for (let i = 0; i < pasted.length; i++) {
        newDigits[i] = pasted[i]
      }
      setDigits(newDigits)

      if (pasted.length === 6) {
        handleVerify(pasted)
      } else {
        inputsRef.current[pasted.length]?.focus()
      }
    }
  }

  const handleVerify = async (code: string) => {
    if (!state.requestId) return

    try {
      await verifyOtp.mutateAsync({
        request_id: state.requestId,
        code,
        device_fingerprint: navigator.userAgent,
      })

      // Success — redirect to home
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'الكود مش صحيح')
      // Clear inputs and shake
      setDigits(Array(6).fill(''))
      inputsRef.current[0]?.focus()
    }
  }

  const handleResend = async () => {
    if (!state.phone) return
    // Re-send OTP by going back to login and resubmitting
    navigate('/login', { replace: true })
  }

  return (
    <div className="min-h-screen bg-bg-primary flex flex-col items-center justify-center p-6">
      <div className="w-full max-w-md">
        {/* Back */}
        <button
          onClick={() => navigate('/login')}
          className="flex items-center gap-1 text-text-secondary hover:text-text-primary mb-8"
        >
          <span className="material-symbols-rounded">arrow_forward</span>
          <span className="text-body">رجوع</span>
        </button>

        {/* Title */}
        <div className="text-center mb-8">
          <h1 className="text-h1 font-bold text-text-primary">أدخل الكود</h1>
          <p className="text-body text-text-secondary mt-2">
            بعتناك كود على
          </p>
          <p className="text-body font-semibold text-text-primary mt-1" dir="ltr">
            +20 {state.phone}
          </p>
          <Link to="/login" className="text-caption text-primary hover:underline mt-1 inline-block">
            تغيير الرقم
          </Link>
        </div>

        {/* OTP inputs */}
        <div
          className="flex justify-center gap-2 mb-6"
          dir="ltr"
          onPaste={handlePaste}
        >
          {digits.map((digit, index) => (
            <input
              key={index}
              ref={(el) => {
                inputsRef.current[index] = el
              }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={digit}
              onChange={(e) => handleChange(index, e.target.value)}
              onKeyDown={(e) => handleKeyDown(index, e)}
              className={`w-12 h-14 text-center text-h2 font-bold rounded-lg border-2 transition-colors
                ${error
                  ? 'border-error bg-error/5'
                  : digit
                    ? 'border-primary bg-primary/5'
                    : 'border-border bg-surface'
                }
                focus:outline-none focus:ring-2 focus:ring-primary/20
              `}
              autoFocus={index === 0}
            />
          ))}
        </div>

        {/* Error */}
        {error && (
          <p className="text-center text-body-sm text-error mb-4" role="alert">
            {error}
          </p>
        )}

        {/* Timer */}
        <div className="text-center mb-6">
          {isRunning && seconds > 0 ? (
            <p className="text-caption text-text-tertiary">
              <span className="material-symbols-rounded text-sm align-middle mr-1">timer</span>
              {formatted} متبقي
            </p>
          ) : (
            <div className="space-y-1">
              <p className="text-caption text-text-tertiary">ماوصلكش الكود؟</p>
              <button
                onClick={handleResend}
                className="text-body-sm text-primary font-semibold hover:underline"
              >
                إعادة إرسال
              </button>
            </div>
          )}
        </div>

        {/* Submit (optional — auto-submits when complete) */}
        <Button
          fullWidth
          size="lg"
          isLoading={verifyOtp.isPending}
          onClick={() => handleVerify(digits.join(''))}
          disabled={digits.some((d) => !d)}
        >
          تأكيد
        </Button>
      </div>
    </div>
  )
}
