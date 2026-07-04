// Phone login page — customer enters phone number to receive OTP

import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useSendOtp } from '@food-platform/hooks'
import { Button, Input } from '@food-platform/ui'
import { isPhoneValid } from '@food-platform/utils'
import type { UserRole } from '@food-platform/types'

export function LoginPage() {
  const navigate = useNavigate()
  const sendOtp = useSendOtp()

  const [phone, setPhone] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!phone) {
      setError('أدخل رقم الموبايل')
      return
    }

    if (!isPhoneValid(phone)) {
      setError('رقم الموبايل مش صحيح')
      return
    }

    try {
      const result = await sendOtp.mutateAsync({
        phone,
        role: 'customer' as UserRole,
      })

      // Navigate to OTP page with request ID
      navigate('/otp', {
        state: {
          requestId: result.request_id,
          phone,
          expiresIn: result.expires_in,
        },
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'حصل خطأ، حاول تاني')
    }
  }

  return (
    <div className="min-h-screen bg-bg-primary flex flex-col items-center justify-center p-6">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-primary rounded-2xl mb-4">
            <span className="material-symbols-rounded text-white text-3xl">restaurant</span>
          </div>
          <h1 className="text-h1 font-bold text-text-primary">أهلاً! 👋</h1>
          <p className="text-body text-text-secondary mt-2">
            سجّل رقم موبايلك عشان نبعتلك كود تفعيل
          </p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="رقم الموبايل"
            type="tel"
            value={phone}
            onChange={(e) => {
              setPhone(e.target.value)
              setError('')
            }}
            placeholder="01 2345 6789"
            error={error}
            startIcon={<span className="text-body">🇪🇬 +20</span>}
            inputMode="numeric"
            maxLength={11}
            autoFocus
          />

          <Button
            type="submit"
            fullWidth
            size="lg"
            isLoading={sendOtp.isPending}
          >
            أرسل الكود
          </Button>
        </form>

        {/* Terms */}
        <p className="text-caption text-text-tertiary text-center mt-6">
          بتسجيلك، أنت موافق على{' '}
          <Link to="/terms" className="text-primary hover:underline">
            الشروط والأحكام
          </Link>{' '}
          و{' '}
          <Link to="/privacy" className="text-primary hover:underline">
            سياسة الخصوصية
          </Link>
        </p>
      </div>
    </div>
  )
}
