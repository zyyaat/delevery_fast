// Checkout page — final review and place order with payment

import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { AppLayout } from '../../components/AppLayout'
import { Button, Badge } from '@food-platform/ui'
import { formatEGP } from '@food-platform/utils'

type PaymentMethod = 'vodafone_cash' | 'instapay' | 'card' | 'cod'

const PAYMENT_OPTIONS: { method: PaymentMethod; label: string; icon: string; description: string; fee?: number }[] = [
  { method: 'vodafone_cash', label: 'Vodafone Cash', icon: '💚', description: 'رصيدك: EGP 1,250' },
  { method: 'instapay', label: 'InstaPay', icon: '🟣', description: 'تحويل بنكي فوري' },
  { method: 'card', label: 'بطاقة بنكية', icon: '💳', description: '**** 4521' },
  { method: 'cod', label: 'الدفع عند الاستلام', icon: '💵', description: 'ادفع كاش للمندوب', fee: 5 },
]

// Mock totals (will come from cart store)
const SUBTOTAL = 550
const DELIVERY_FEE = 25
const SERVICE_FEE = SUBTOTAL * 0.05
const VAT = (SUBTOTAL + DELIVERY_FEE + SERVICE_FEE) * 0.14
const DISCOUNT = 50
const TOTAL = SUBTOTAL + DELIVERY_FEE + SERVICE_FEE + VAT - DISCOUNT

export function CheckoutPage() {
  const navigate = useNavigate()
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('vodafone_cash')
  const [notes, setNotes] = useState('')
  const [isPlacing, setIsPlacing] = useState(false)

  const codFee = paymentMethod === 'cod' ? 5 : 0
  const finalTotal = TOTAL + codFee

  const handlePlaceOrder = async () => {
    setIsPlacing(true)

    // Simulate API call (in production: create order → charge payment → redirect)
    await new Promise((resolve) => setTimeout(resolve, 1500))

    // Navigate to order tracking with mock order ID
    navigate('/orders/mock-order-id', {
      state: {
        orderConfirmed: true,
        orderId: 'mock-order-id',
        total: finalTotal,
        eta: '8:35 PM',
      },
    })
  }

  return (
    <AppLayout>
      <div className="max-w-2xl mx-auto px-4 py-4 pb-28">
        <h1 className="text-h2 font-bold text-text-primary mb-4">إتمام الطلب</h1>

        {/* Address */}
        <Section title="📍 عنوان التوصيل">
          <div className="bg-surface rounded-lg border border-border p-4">
            <p className="text-body font-medium text-text-primary">
              الزمالك، 26 يوليو، شقة 5
            </p>
            <p className="text-caption text-text-secondary mt-1">علامة: باب أزرق</p>
            <button className="text-caption text-primary font-semibold mt-2 hover:underline">
              تغيير
            </button>
          </div>
        </Section>

        {/* Time */}
        <Section title="⏰ وقت التوصيل">
          <div className="bg-surface rounded-lg border border-border p-4 space-y-2">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="time"
                defaultChecked
                className="w-5 h-5 accent-primary"
              />
              <div>
                <p className="text-body-sm font-medium text-text-primary">في أسرع وقت</p>
                <p className="text-caption text-text-secondary">35-45 دقيقة</p>
              </div>
            </label>
            <label className="flex items-center gap-3 cursor-pointer">
              <input type="radio" name="time" className="w-5 h-5 accent-primary" />
              <div>
                <p className="text-body-sm font-medium text-text-primary">مجدول لوقت لاحق</p>
                <p className="text-caption text-text-secondary">اختر التاريخ والوقت</p>
              </div>
            </label>
          </div>
        </Section>

        {/* Payment Method */}
        <Section title="💳 طريقة الدفع">
          <div className="bg-surface rounded-lg border border-border p-4 space-y-3">
            {PAYMENT_OPTIONS.map((option) => (
              <label
                key={option.method}
                className={`flex items-center gap-3 p-3 rounded-lg border-2 cursor-pointer transition-colors
                  ${paymentMethod === option.method
                    ? 'border-primary bg-primary/5'
                    : 'border-border hover:border-border-strong'
                  }`}
              >
                <input
                  type="radio"
                  name="payment"
                  checked={paymentMethod === option.method}
                  onChange={() => setPaymentMethod(option.method)}
                  className="w-5 h-5 accent-primary"
                />
                <span className="text-2xl">{option.icon}</span>
                <div className="flex-1">
                  <p className="text-body-sm font-semibold text-text-primary">{option.label}</p>
                  <p className="text-caption text-text-secondary">{option.description}</p>
                </div>
                {option.fee && (
                  <Badge variant="warning">+{formatEGP(option.fee)}</Badge>
                )}
              </label>
            ))}
          </div>
        </Section>

        {/* Coupon */}
        <Section title="🎁 كوبون">
          <div className="bg-surface rounded-lg border border-border p-4 flex items-center justify-between">
            <span className="text-body-sm text-success font-medium">✅ WELCOME50 مطبق</span>
            <span className="text-body-sm text-success">-{formatEGP(DISCOUNT)}</span>
          </div>
        </Section>

        {/* Notes */}
        <Section title="📝 ملاحظات للمطعم">
          <input
            type="text"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            placeholder="مثلاً: صلصة إضافية لو سمحت"
            maxLength={200}
            className="w-full h-12 px-4 rounded-lg border border-border bg-surface text-body placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </Section>

        {/* Phone */}
        <Section title="📞 رقم التواصل">
          <div className="bg-surface rounded-lg border border-border p-4 flex items-center justify-between">
            <span className="text-body-sm text-text-primary" dir="ltr">010 1234 5678</span>
            <button className="text-caption text-primary font-semibold hover:underline">تغيير</button>
          </div>
        </Section>

        {/* Total */}
        <div className="bg-surface rounded-lg border border-border p-4 mt-4">
          <div className="flex items-center justify-between">
            <span className="text-h3 font-bold text-text-primary">الإجمالي</span>
            <span className="text-h3 font-bold text-primary">{formatEGP(finalTotal)}</span>
          </div>
        </div>

        {/* Place order button (sticky) */}
        <div className="fixed bottom-16 md:bottom-0 left-0 right-0 bg-surface border-t border-border p-4 z-sticky">
          <div className="max-w-2xl mx-auto">
            <Button
              fullWidth
              size="lg"
              isLoading={isPlacing}
              onClick={handlePlaceOrder}
            >
              {isPlacing ? 'جاري تأكيد الطلب...' : `🔒 تأكيد الطلب — ${formatEGP(finalTotal)}`}
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mb-4">
      <h2 className="text-body font-semibold text-text-primary mb-2">{title}</h2>
      {children}
    </div>
  )
}
